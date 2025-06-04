package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	jwtAccessSecret  string
	jwtRefreshSecret string
}

type Claims struct {
	Email       string `json:"email"`
	UserID      int    `json:"user_id,omitempty"`
	IsActivated bool   `json:"is_activated,omitempty"`
	jwt.RegisteredClaims
}

func New(jwtAccessSecret, jwtRefreshSecret string) *TokenService {
	return &TokenService{
		jwtAccessSecret:  jwtAccessSecret,
		jwtRefreshSecret: jwtRefreshSecret,
	}
}

func (t *TokenService) GenerateTokens(email string, userID int, isActivated bool) (map[string]string, error) {

	accessSecret := []byte(t.jwtAccessSecret)
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Email:       email,
		UserID:      userID,
		IsActivated: isActivated,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}).SignedString(accessSecret)
	if err != nil {
		return nil, err
	}

	refreshSecret := []byte(t.jwtRefreshSecret)
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Email:       email,
		UserID:      userID,
		IsActivated: isActivated,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}).SignedString(refreshSecret)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}
