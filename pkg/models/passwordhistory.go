package models

type PasswordHistory struct {
	ID                int    `json:"id"`
	Username          string `json:"username"`
	PreviousPassword1 string `json:"previous_password_1"`
	PreviousPassword2 string `json:"previous_password_2"`
	PreviousPassword3 string `json:"previous_password_3"`
}

func NewPasswordHistory(username string, previous1, previous2, previous3 string) PasswordHistory {
	return PasswordHistory{
		Username:          username,
		PreviousPassword1: previous1,
		PreviousPassword2: previous2,
		PreviousPassword3: previous3,
	}
}

func (PasswordHistory) TableName() string {
	return "password_history"
}
