package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/middleware/security"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type ShowController struct {
	showService services.ShowService
}

func NewShowController(showService services.ShowService) *ShowController {
	return &ShowController{
		showService: showService,
	}
}

func (sh *ShowController) GetShows(c *gin.Context) {
	requestID := utils.GetRequestID(c)
	
	dateStr := c.Query("date")
	requestDate, err := utils.GetDateFromDateStringDefaultToday(dateStr)
	
	if err != nil {
		utils.HandleErrorResponse(c, utils.NewBadRequestError("INVALID_DATE", "Invalid date format. Use YYYY-MM-DD", err), requestID)
		return
	}

	claims, err := security.GetTokenClaims(c)
	if err != nil {
		utils.HandleErrorResponse(c, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify user credentials", err), requestID)
		return
	}

	role, _ := claims["role"].(string)
	username, _ := claims["username"].(string)
	if role != "admin" && role != "staff" {
		today := time.Now().Truncate(24 * time.Hour)
		maxAllowedDate := today.AddDate(0, 0, 6)

		if requestDate.Before(today) || requestDate.After(maxAllowedDate) {
			log.Info().
				Str("username", username).
				Time("requestedDate", requestDate).
				Msg("Customer attempted to access shows outside allowed date range")

			utils.HandleErrorResponse(c,
				utils.NewBadRequestError(
					"DATE_OUT_OF_RANGE",
					"Customers can only view shows from today to the next 6 days",
					nil),
				requestID)
			return
		}
	}

	shows, err := sh.showService.GetShows(c.Request.Context(), requestDate)
	if err != nil {
		log.Error().Err(err).Str("date", dateStr).Msg("Failed to get shows")
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	var showResponses []response.ShowResponse
	for _, show := range shows {
		movie, err := sh.showService.GetMovieById(c.Request.Context(), show.MovieId)
		if err != nil {
			log.Error().Err(err).Str("movieId", show.MovieId).Msg("Failed to get movie for show")
			utils.HandleErrorResponse(c, err, requestID)
			return
		}

		availableSeats := sh.showService.AvailableSeats(c.Request.Context(), show.Id)
		showResponse := response.NewShowResponse(*movie, show.Slot, show, availableSeats)
		showResponses = append(showResponses, *showResponse)
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message:   "Shows retrieved successfully",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data:      showResponses,
	})
}

func (sh *ShowController) GetMovies(c *gin.Context) {
	requestID := utils.GetRequestID(c)
	movies, err := sh.showService.GetMovies(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get movies")
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	var movieResponses []response.MovieResponse
	for _, movie := range movies {
		movieResponses = append(movieResponses, response.NewMovieResponse(movie))
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message:   "Movies retrieved successfully",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data:      movieResponses,
	})
}

func (sh *ShowController) CreateShow(c *gin.Context) {
	requestID := utils.GetRequestID(c)
	var showRequest request.ShowRequest
	if err := c.ShouldBindJSON(&showRequest); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			utils.HandleErrorResponse(c, utils.NewValidationError(validationErrs), requestID)
			return
		}
		utils.HandleErrorResponse(c, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request data", err), requestID)
		return
	}

	show, err := sh.showService.CreateShow(c.Request.Context(), showRequest)
	if err != nil {
		log.Error().Err(err).Interface("request", showRequest).Msg("Failed to create show")
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse{
		Message:   "Show created successfully",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data: response.NewShowConfirmationResponse(
			show.Id,
			show.MovieId,
			show.Slot,
			show.Date.Format("2006-01-02"),
			show.Cost,
		),
	})
}
