package models

type Slot struct {
    Id        int    `json:"id"`
    Name      string `json:"name"`
    StartTime string `json:"startTime"`
    EndTime   string `json:"endTime"`
}

