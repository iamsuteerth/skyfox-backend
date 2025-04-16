package controllers

import (
	"encoding/base64"
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

func (sk *SkyCustomerController) Signup(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)
	var req request.SignupRequest

	if err := sk.parseAndValidateRequest(ctx, &req); err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	if err := sk.skyCustomerService.ValidateUserDetails(ctx.Request.Context(), req.Username, req.Email, req.PhoneNumber); err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	if err := sk.securityQuestionService.ValidateSecurityQuestionExists(ctx.Request.Context(), req.SecurityQuestionID); err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewInternalServerError("PASSWORD_HASH_ERROR", "Failed to hash password", err), requestID)
		return
	}

	imageBytes, err := base64.StdEncoding.DecodeString(req.ProfileImg)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_IMAGE", "Invalid base64 image data", err), requestID)
		return
	}

	user := models.NewUser(req.Username, hashedPassword, "customer")
	passwordHistory := models.NewPasswordHistory(req.Username, hashedPassword, "", "")
	skyCustomer := models.NewSkyCustomer(req.Name, req.Username, req.PhoneNumber, req.Email, imageBytes, 0, "")

	if err := sk.skyCustomerService.CreateCustomer(
		ctx.Request.Context(),
		&skyCustomer,
		&user,
		&passwordHistory,
		req.SecurityQuestionID,
		req.SecurityAnswer,
		imageBytes,
		req.ProfileImgSHA,
	); err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	response := response.SignupResponse{
		Username: req.Username,
		Name:     req.Name,
	}

	utils.SendCreatedResponse(ctx, "User registered successfully", requestID, response)
}

func (sk *SkyCustomerController) GetCustomerProfile(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify credentials", err), requestID)
		return
	}

	tokenUsername, ok := claims["username"].(string)
	if !ok || tokenUsername == "" {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("INVALID_TOKEN", "Invalid token claims", nil), requestID)
		return
	}

	profile, err := sk.skyCustomerService.GetCustomerProfile(ctx.Request.Context(), tokenUsername)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	utils.SendOKResponse(ctx, "Customer profile retrieved successfully", requestID, profile)
}

func (sk *SkyCustomerController) GetProfileImagePresignedURL(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify credentials", err), requestID)
		return
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("INVALID_TOKEN", "Invalid token claims", nil), requestID)
		return
	}

	presignedURL, expiresAt, err := sk.skyCustomerService.GetProfileImagePresignedURL(ctx.Request.Context(), username)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	response := response.ProfileImageResponse{
		PresignedURL: presignedURL,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
	}

	utils.SendOKResponse(ctx, "Presigned URL generated successfully", requestID, response)
}

func (sc *SkyCustomerController) UpdateCustomerProfile(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify credentials", err), requestID)
		return
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("INVALID_TOKEN", "Invalid token claims", nil), requestID)
		return
	}

	var updateRequest request.UpdateCustomerProfileRequest
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			utils.HandleErrorResponse(ctx, utils.NewValidationError(validationErrs), requestID)
			return
		}
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request data", err), requestID)
		return
	}

	customer, err := sc.skyCustomerService.GetCustomerProfile(ctx.Request.Context(), username)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	verifyResponse, err := sc.securityQuestionService.VerifySecurityAnswer(
		ctx.Request.Context(),
		customer.Email,
		updateRequest.SecurityAnswer,
	)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	if verifyResponse == nil || !verifyResponse.ValidAnswer {
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_SECURITY_ANSWER", "The security answer provided is incorrect", nil), requestID)
		return
	}

	updatedProfile, err := sc.skyCustomerService.UpdateCustomerProfile(
		ctx.Request.Context(),
		username,
		&updateRequest,
	)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	utils.SendOKResponse(ctx, "Profile updated successfully", requestID, updatedProfile)
}

func (sc *SkyCustomerController) UpdateProfileImage(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify credentials", err), requestID)
		return
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("INVALID_TOKEN", "Invalid token claims", nil), requestID)
		return
	}

	var updateRequest request.UpdateProfileImageRequest
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			utils.HandleErrorResponse(ctx, utils.NewValidationError(validationErrs), requestID)
			return
		}
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request data", err), requestID)
		return
	}

	customer, err := sc.skyCustomerService.GetCustomerProfile(ctx.Request.Context(), username)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	verifyResponse, err := sc.securityQuestionService.VerifySecurityAnswer(
		ctx.Request.Context(),
		customer.Email,
		updateRequest.SecurityAnswer,
	)

	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	if verifyResponse == nil || !verifyResponse.ValidAnswer {
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_SECURITY_ANSWER", "The security answer provided is incorrect", nil), requestID)
		return
	}

	imageBytes, err := base64.StdEncoding.DecodeString(updateRequest.ProfileImg)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_IMAGE", "Invalid base64 image data", err), requestID)
		return
	}

	err = sc.skyCustomerService.UpdateProfileImage(
		ctx.Request.Context(),
		username,
		imageBytes,
		updateRequest.ProfileImgSHA,
	)

	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	utils.SendOKResponse(ctx, "Profile image updated successfully", requestID, nil)
}

func (sk *SkyCustomerController) parseAndValidateRequest(ctx *gin.Context, req *request.SignupRequest) error {
	contentType := ctx.Request.Header.Get("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
			return utils.NewBadRequestError("INVALID_FORM", "Failed to parse multipart form", err)
		}

		req.Name = ctx.Request.FormValue("name")
		req.Username = ctx.Request.FormValue("username")
		req.Email = ctx.Request.FormValue("email")
		req.PhoneNumber = ctx.Request.FormValue("number")
		req.Password = ctx.Request.FormValue("password")
		req.SecurityAnswer = ctx.Request.FormValue("security_answer")
		req.ProfileImgSHA = ctx.Request.FormValue("profile_img_sha")

		securityQuestionIDStr := ctx.Request.FormValue("security_question_id")
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

		file, header, err := ctx.Request.FormFile("profile_img")
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
			req.ProfileImg = ctx.Request.FormValue("profile_img")
			if req.ProfileImg == "" {
				return utils.NewBadRequestError("MISSING_PROFILE_IMAGE", "Profile image is required", nil)
			}
		}

		if req.ProfileImgSHA == "" {
			return utils.NewBadRequestError("MISSING_IMAGE_HASH", "Profile image hash is required", nil)
		}
	} else {
		if err := ctx.ShouldBindJSON(req); err != nil {
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
