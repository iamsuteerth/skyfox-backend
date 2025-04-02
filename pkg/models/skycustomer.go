package models

type SkyCustomer struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Username           string `json:"username"`
	Number             string `json:"number"`
	Email              string `json:"email"`
	ProfileImg         []byte `json:"profile_img,omitempty"`
	SecurityQuestionID int    `json:"security_question_id"`
	SecurityAnswerHash string `json:"-"`
}

func NewSkyCustomer(name, username, number, email string, profileImg []byte, securityQuestionID int, securityAnswerHash string) SkyCustomer {
	return SkyCustomer{
		Name:               name,
		Username:           username,
		Number:             number,
		Email:              email,
		ProfileImg:         profileImg,
		SecurityQuestionID: securityQuestionID,
		SecurityAnswerHash: securityAnswerHash,
	}
}

func (SkyCustomer) TableName() string {
	return "customertable"
}
