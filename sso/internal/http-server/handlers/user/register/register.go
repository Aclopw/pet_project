package register

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sso/internal/domain/models"
	"sso/internal/lib/logger/sl"
	"sso/internal/lib/random"
	"time"

	resp "sso/internal/lib/api/respones"
	"sso/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=64"`
}

type Response struct {
	resp.Response
	models.UserDTO `json:"user,omitempty"`
	AccessToken    string `json:"access_token,omitempty"`
	RefreshToken   string `json:"refresh_token,omitempty"`
}

type SendVerificationEmail interface {
	Send(email, activationLink string) error
}

type TokenService interface {
	GenerateTokens(email string, userID int, isActivated bool) (map[string]string, error)
}

// TokenService abstracts token generation logic
type Saver interface {
	SaveUser(user *models.User) (int, error)
	SaveToken(userID int, token string) error
}

func New(log *slog.Logger, userSaver Saver, mailService SendVerificationEmail, tokenService TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.register.New"

		// Add operation and request ID to logger context.
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Decode the request body into the Request struct.
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Validate the request struct.
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		// Hash the user's password.
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		activationLink := random.NewRandomString(27)

		user := models.User{Email: req.Email, Password: string(hashedPassword), TgUserID: "", ActivationLink: activationLink}

		// Save the new user to the storage.
		userID, err := userSaver.SaveUser(&user)
		if err != nil {
			if errors.Is(err, storage.ErrUserAlreadyExists) {
				log.Error("user already exists", sl.Err(err))

				render.JSON(w, r, resp.Error("user already exists"))

				return
			}
			log.Error("failed to save user", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		// Send verification email asynchronously.
		go func() {
			if err := mailService.Send(req.Email, activationLink); err != nil {
				log.Error("failed to send verification email", sl.Err(err))
			}
		}()

		// Generate access and refresh tokens for the user.
		tokens, err := tokenService.GenerateTokens(req.Email, userID, false)
		if err != nil {
			log.Error("failed to generate tokens", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		// Set the refresh token as a cookie.
		cookie := http.Cookie{Name: "resfreshToken", Value: tokens["refresh_token"], Expires: cookieExpiration(), HttpOnly: true}
		http.SetCookie(w, &cookie)

		// Save the refresh token in storage.
		if err := userSaver.SaveToken(userID, tokens["refresh_token"]); err != nil {
			log.Error("failed to save token", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		responesOk(w, r, models.UserDTO{
			Email:       req.Email,
			UserID:      userID,
			IsActivated: false,
		}, tokens)

		log.Info("user registered successfully", slog.String("email", req.Email), slog.Int("user_id", userID))
	}
}

func responesOk(w http.ResponseWriter, r *http.Request, user models.UserDTO, tokens map[string]string) {
	resp := Response{
		Response:     resp.OK(),
		UserDTO:      user,
		AccessToken:  tokens["access_token"],
		RefreshToken: tokens["refresh_token"],
	}

	render.JSON(w, r, resp)
}

// cookieExpiration returns the expiration time for the cookie (30 days from now)
func cookieExpiration() time.Time {
	return time.Now().Add(60 * 24 * 30 * time.Minute)
}
