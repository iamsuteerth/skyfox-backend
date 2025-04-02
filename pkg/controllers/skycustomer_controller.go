package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type SkyCustomerController struct {
	userService        services.UserService
	skyCustomerService services.SkyCustomerService
}

func NewSkyCustomerController(userService services.UserService, skyCustomerService services.SkyCustomerService) *SkyCustomerController {
	return &SkyCustomerController{
		userService:        userService,
		skyCustomerService: skyCustomerService,
	}
}

type SignupRequest struct {
	Name        string `json:"name" binding:"required,customName"`
	Username    string `json:"username" binding:"required,customUsername"`
	Password    string `json:"password" binding:"required,min=8,max=32"`
	PhoneNumber string `json:"number" binding:"required,customPhone"`
	Email       string `json:"email" binding:"required,email"`
	ProfileImg  []byte `json:"profile_img"`
}

func (sk *SkyCustomerController) Signup(c *gin.Context) {
	requestID := utils.GetRequestID(c)
	var req SignupRequest

	if err := sk.parseAndValidateRequest(c, &req); err != nil {
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	if err := sk.skyCustomerService.ValidateUserDetails(c.Request.Context(), req.Username, req.Email, req.PhoneNumber); err != nil {
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.HandleErrorResponse(c, utils.NewInternalServerError("PASSWORD_HASH_ERROR", "Failed to hash password", err), requestID)
		return
	}

	user := models.NewUser(req.Username, hashedPassword, "customer")
	passwordHistory := models.NewPasswordHistory(req.Username, hashedPassword, "", "")
	skyCustomer := models.NewSkyCustomer(req.Name, req.Username, req.PhoneNumber, req.Email, req.ProfileImg)

	if err := sk.skyCustomerService.CreateCustomerWithUser(c.Request.Context(), user, passwordHistory, skyCustomer); err != nil {
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse{
		Message:   "User registered successfully",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data: gin.H{
			"username": req.Username,
			"name":     req.Name,
		},
	})
}

func (sk *SkyCustomerController) parseAndValidateRequest(c *gin.Context, req *SignupRequest) error {
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

		if req.Name == "" || req.Username == "" || req.Email == "" || req.PhoneNumber == "" || req.Password == "" {
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
			req.ProfileImg = make([]byte, header.Size)
			if _, err := file.Read(req.ProfileImg); err != nil {
				return utils.NewInternalServerError("FILE_READ_ERROR", "Failed to read file", err)
			}
		}
	} else {
		if err := c.ShouldBindJSON(req); err != nil {
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				return utils.NewValidationError(validationErrs)
			}
			return utils.NewBadRequestError("INVALID_REQUEST", "Invalid request data", err)
		}
	}

	req.Name = utils.ToCamelCase(strings.TrimSpace(req.Name))

	return nil
}
