package request

import "github.com/iamsuteerth/skyfox-backend/pkg/models"

type ShowRequest struct {
    MovieId string      `json:"movieId"`
    Date    string      `json:"date" binding:"required,datetime=2006-01-02"`
    Slot    models.Slot `json:"slot"`
    SlotId  int         `json:"slotId"`
    Cost    float64     `json:"cost"`
}

