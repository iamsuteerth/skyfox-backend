package controllers

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/middleware/security"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type SkyCustomerController struct {
	userService             services.UserService
	skyCustomerService      services.SkyCustomerService
	securityQuestionService services.SecurityQuestionService
}

func NewSkyCustomerController(
	userService services.UserService,
	skyCustomerService services.SkyCustomerService,
	securityQuestionService services.SecurityQuestionService,
) *SkyCustomerController {
	return &SkyCustomerController{
		userService:             userService,
		skyCustomerService:      skyCustomerService,
		securityQuestionService: securityQuestionService,
	}
}

func (sk *SkyCustomerController) Signup(c *gin.Context) {
	requestID := utils.GetRequestID(c)
	var req request.SignupRequest

	if err := sk.parseAndValidateRequest(c, &req); err != nil {
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	if err := sk.skyCustomerService.ValidateUserDetails(c.Request.Context(), req.Username, req.Email, req.PhoneNumber); err != nil {
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	if err := sk.securityQuestionService.ValidateSecurityQuestionExists(c.Request.Context(), req.SecurityQuestionID); err != nil {
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.HandleErrorResponse(c, utils.NewInternalServerError("PASSWORD_HASH_ERROR", "Failed to hash password", err), requestID)
		return
	}

	imageBytes, err := base64.StdEncoding.DecodeString(req.ProfileImg)
	if err != nil {
		utils.HandleErrorResponse(c, utils.NewBadRequestError("INVALID_IMAGE", "Invalid base64 image data", err), requestID)
		return
	}

	user := models.NewUser(req.Username, hashedPassword, "customer")
	passwordHistory := models.NewPasswordHistory(req.Username, hashedPassword, "", "")
	skyCustomer := models.NewSkyCustomer(req.Name, req.Username, req.PhoneNumber, req.Email, imageBytes, 0, "")

	if err := sk.skyCustomerService.CreateCustomer(
		c.Request.Context(),
		&skyCustomer,
		&user,
		&passwordHistory,
		req.SecurityQuestionID,
		req.SecurityAnswer,
		imageBytes,
		req.ProfileImgSHA,
	); err != nil {
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse{
		Message:   "User registered successfully",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data: response.SignupResponse{
			Username: req.Username,
			Name:     req.Name,
		},
	})
}

func (sk *SkyCustomerController) parseAndValidateRequest(c *gin.Context, req *request.SignupRequest) error {
	contentType := c.Request.Header.Get("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			return utils.NewBadRequestError("INVALID_FORM", "Failed to parse multipart form", err)
		}

		req.Name = c.Request.FormValue("name")
		req.Username = c.Request.FormValue("username")
		req.Email = c.Request.FormValue("email")
		req.PhoneNumber = c.Request.FormValue("number")
		req.Password = c.Request.FormValue("password")
		req.SecurityAnswer = c.Request.FormValue("security_answer")
		req.ProfileImgSHA = c.Request.FormValue("profile_img_sha")

		securityQuestionIDStr := c.Request.FormValue("security_question_id")
		if securityQuestionIDStr != "" {
			securityQuestionID, err := strconv.Atoi(securityQuestionIDStr)
			if err != nil {
				return utils.NewBadRequestError("INVALID_SECURITY_QUESTION", "Security question ID must be a valid number", err)
			}
			req.SecurityQuestionID = securityQuestionID
		}

		if req.Name == "" || req.Username == "" || req.Email == "" || req.PhoneNumber == "" ||
			req.Password == "" || req.SecurityQuestionID == 0 || req.SecurityAnswer == "" {
			return utils.NewBadRequestError("MISSING_FIELDS", "All fields are required", nil)
		}

		tempReq := *req
		if err := binding.Validator.ValidateStruct(tempReq); err != nil {
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				return utils.NewValidationError(validationErrs)
			}
		}

		file, header, err := c.Request.FormFile("profile_img")
		if err == nil && header != nil {
			defer file.Close()
			if header.Size > 5<<20 {
				return utils.NewBadRequestError("FILE_TOO_LARGE", "File size exceeds 5MB limit", nil)
			}

			imageBytes := make([]byte, header.Size)
			if _, err := file.Read(imageBytes); err != nil {
				return utils.NewInternalServerError("FILE_READ_ERROR", "Failed to read file", err)
			}

			req.ProfileImg = base64.StdEncoding.EncodeToString(imageBytes)
		} else {
			req.ProfileImg = c.Request.FormValue("profile_img")
			if req.ProfileImg == "" {
				return utils.NewBadRequestError("MISSING_PROFILE_IMAGE", "Profile image is required", nil)
			}
		}

		if req.ProfileImgSHA == "" {
			return utils.NewBadRequestError("MISSING_IMAGE_HASH", "Profile image hash is required", nil)
		}
	} else {
		if err := c.ShouldBindJSON(req); err != nil {
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				return utils.NewValidationError(validationErrs)
			}
			return utils.NewBadRequestError("INVALID_REQUEST", "Invalid request data", err)
		}

		if req.SecurityQuestionID == 0 || req.SecurityAnswer == "" {
			return utils.NewBadRequestError("MISSING_SECURITY_FIELDS", "Security question and answer are required", nil)
		}

		if req.ProfileImg == "" || req.ProfileImgSHA == "" {
			return utils.NewBadRequestError("MISSING_PROFILE_IMAGE", "Profile image and hash are required", nil)
		}
	}

	req.Name = utils.ToCamelCase(strings.TrimSpace(req.Name))
	req.SecurityAnswer = strings.TrimSpace(req.SecurityAnswer)

	return nil
}

func (sk *SkyCustomerController) GetProfileImagePresignedURL(c *gin.Context) {
	requestID := utils.GetRequestID(c)

	claims, err := security.GetTokenClaims(c)
	if err != nil {
		utils.HandleErrorResponse(c, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify credentials", err), requestID)
		return
	}

	tokenUsername, ok := claims["username"].(string)
	if !ok || tokenUsername == "" {
		utils.HandleErrorResponse(c, utils.NewUnauthorizedError("INVALID_TOKEN", "Invalid token claims", nil), requestID)
		return
	}

	username := tokenUsername

	presignedURL, expiresAt, err := sk.skyCustomerService.GetProfileImagePresignedURL(c.Request.Context(), username)
	if err != nil {
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message:   "Presigned URL generated successfully",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data: response.ProfileImageResponse{
			PresignedURL: presignedURL,
			ExpiresAt:    expiresAt.Format(time.RFC3339),
		},
	})
}
