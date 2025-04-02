package models

type Staff struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	Name          string `json:"name"`
	CounterNumber int    `json:"counter_number"`
}

func NewStaff(username string, name string, counterNumber int) Staff {
	return Staff{
		Username:      username,
		Name:          name,
		CounterNumber: counterNumber,
	}
}

func (Staff) TableName() string {
	return "stafftable"
}
