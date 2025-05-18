package response

import (
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
)

type ShowResponse struct {
	Movie          models.Movie `json:"movie"`
	Slot           models.Slot  `json:"slot"`
	Id             int          `json:"id"`
	Date           time.Time    `json:"date"`
	Cost           float64      `json:"cost"`
	AvailableSeats int          `json:"availableseats"`
}

func NewShowResponse(movie models.Movie, slot models.Slot, show models.Show, availableSeats int) *ShowResponse {
	showCost, _ := show.Cost.Float64()
	return &ShowResponse{
		Movie:          movie,
		Slot:           slot,
		Id:             show.Id,
		Date:           show.Date,
		Cost:           showCost,
		AvailableSeats: availableSeats,
	}
}

type ShowConfirmationResponse struct {
	Id      int         `json:"id"`
	MovieId string      `json:"movie"`
	Slot    models.Slot `json:"slot"`
	Date    string      `json:"date"`
	Cost    float64     `json:"cost"`
}

func NewShowConfirmationResponse(id int, movieId string, slot models.Slot, date string, cost float64) *ShowConfirmationResponse {
	return &ShowConfirmationResponse{
		Id:      id,
		MovieId: movieId,
		Slot:    slot,
		Date:    date,
		Cost:    cost,
	}
}
