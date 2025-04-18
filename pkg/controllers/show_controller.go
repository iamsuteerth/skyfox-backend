package controllers

import (
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

func (sh *ShowController) GetShows(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	dateStr := ctx.Query("date")
	requestDate, err := utils.GetDateFromDateStringDefaultToday(dateStr)

	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_DATE", "Invalid date format. Use YYYY-MM-DD", err), requestID)
		return
	}

	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify user credentials", err), requestID)
		return
	}

	role, _ := claims["role"].(string)
	username, _ := claims["username"].(string)
	if role != "admin" && role != "staff" {
		now := time.Now()

		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		maxAllowedDate := time.Date(today.Year(), today.Month(), today.Day()+6, 23, 59, 59, 0, today.Location())

		requestDateLocal := requestDate.In(today.Location())
		requestDateMidnight := time.Date(requestDateLocal.Year(), requestDateLocal.Month(), requestDateLocal.Day(), 0, 0, 0, 0, today.Location())

		if requestDateMidnight.Before(today) || requestDateMidnight.After(maxAllowedDate) {
			log.Info().
				Str("username", username).
				Time("requestedDate", requestDate).
				Msg("Customer attempted to access shows outside allowed date range")

			utils.HandleErrorResponse(ctx,
				utils.NewBadRequestError(
					"DATE_OUT_OF_RANGE",
					"Customers can only view shows from today to the next 6 days",
					nil),
				requestID)
			return
		}
	}

	shows, err := sh.showService.GetShows(ctx.Request.Context(), requestDate)
	if err != nil {
		log.Error().Err(err).Str("date", dateStr).Msg("Failed to get shows")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	var showResponses []response.ShowResponse
	for _, show := range shows {
		movie, err := sh.showService.GetMovieById(ctx.Request.Context(), show.MovieId)
		if err != nil {
			log.Error().Err(err).Str("movieId", show.MovieId).Msg("Failed to get movie for show")
			utils.HandleErrorResponse(ctx, err, requestID)
			return
		}

		availableSeats := sh.showService.AvailableSeats(ctx.Request.Context(), show.Id)
		showResponse := response.NewShowResponse(*movie, show.Slot, show, availableSeats)
		showResponses = append(showResponses, *showResponse)
	}

	utils.SendOKResponse(ctx, "Shows retrieved successfully", requestID, showResponses)
}

func (sh *ShowController) GetMovies(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)
	movies, err := sh.showService.GetMovies(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get movies")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	var movieResponses []response.MovieResponse
	for _, movie := range movies {
		movieResponses = append(movieResponses, response.NewMovieResponse(movie))
	}

	utils.SendOKResponse(ctx, "Movies retrieved successfully", requestID, movieResponses)
}

func (sh *ShowController) CreateShow(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)
	var showRequest request.ShowRequest
	if err := ctx.ShouldBindJSON(&showRequest); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			utils.HandleErrorResponse(ctx, utils.NewValidationError(validationErrs), requestID)
			return
		}
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request data", err), requestID)
		return
	}

	show, err := sh.showService.CreateShow(ctx.Request.Context(), showRequest)
	if err != nil {
		log.Error().Err(err).Interface("request", showRequest).Msg("Failed to create show")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	response := response.NewShowConfirmationResponse(
		show.Id,
		show.MovieId,
		show.Slot,
		show.Date.Format("2006-01-02"),
		show.Cost,
	)

	utils.SendCreatedResponse(ctx, "Show created successfully", requestID, response)
}
