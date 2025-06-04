package models

type User struct {
	ID             int    `bson:"_id,omitempty"`
	Email          string `bson:"email,omitempty"`
	Password       string `bson:"password,omitempty"`
	IsActive       bool   `bson:"is_active,omitempty"`
	ActivationLink string `bson:"activation_link"`
}
