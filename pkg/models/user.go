package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(username string, password string, role string) User {
	return User{
		Username: username,
		Password: password,
		Role:     role,
	}
}

func (User) TableName() string {
	return "usertable"
}
