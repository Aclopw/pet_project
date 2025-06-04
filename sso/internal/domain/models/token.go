package models

type Token struct {
	ID           int    `bson:"_id,omitempty"`
	UserID       string `bson:"user_id,omitempty"`
	AccessToken  string `bson:"access_token,omitempty"`
	RefreshToken string `bson:"refresh_token,omitempty"`
}
