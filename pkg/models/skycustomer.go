package models

type SkyCustomer struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Number     string `json:"number"`
	Email      string `json:"email"`
	ProfileImg []byte `json:"profile_img"`
}

func NewSkyCustomer(name, username, number, email string, profileImg []byte) SkyCustomer {
	return SkyCustomer{
		Name:       name,
		Username:   username,
		Number:     number,
		Email:      email,
		ProfileImg: profileImg,
	}
}

func (SkyCustomer) TableName() string {
	return "customertable"
}
