package response

import "github.com/iamsuteerth/skyfox-backend/pkg/models"

type SlotResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

func NewSlotResponse(slot models.Slot) SlotResponse {
	return SlotResponse{
		Id:        slot.Id,
		Name:      slot.Name,
		StartTime: slot.StartTime,
		EndTime:   slot.EndTime,
	}
}
