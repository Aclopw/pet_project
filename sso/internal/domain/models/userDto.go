package models

type UserDTO struct {
	Email       string `json:"email"`
	UserID      int    `json:"user_id,omitempty"`
	IsActivated bool   `json:"is_activated,omitempty"`
}
