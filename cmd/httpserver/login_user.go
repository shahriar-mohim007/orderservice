package httpserver

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"orderservice/state"
	utils "orderservice/utils"
	"time"
)

type LoginRequestPayload struct {
	Email    string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponsePayload struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func HandleLogin(app *state.State) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		request := LoginRequestPayload{}
		ctx := req.Context()

		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Invalid JSON",
			})
			_ = ValidDataNotFound.WriteToResponse(w, nil)
			return
		}

		validate := validator.New()

		err = validate.Struct(request)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Invalid payload",
			})
			_ = ValidDataNotFound.WriteToResponse(w, nil)
			return
		}

		user, err := app.Repository.GetUserByEmail(ctx, request.Email)
		if err != nil {
			_ = InvalidEmailPassword.WriteToResponse(w, nil)
			return
		}

		if !utils.CheckPasswordHash(user.Password, request.Password) {
			_ = InvalidEmailPassword.WriteToResponse(w, nil)
			return
		}

		ttl := 24 * time.Hour

		accessToken, err := utils.GenerateJWT(user.ID,
			utils.ScopeAuthentication, app.Config.SecretKey, ttl)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error generating access token",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}

		refreshToken, err := utils.GenerateRefreshToken(user.ID.String(),
			app.Config.SecretKey)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error generating refresh token",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}

		response := LoginResponsePayload{
			TokenType:    "Bearer",
			ExpiresIn:    int64(ttl.Seconds()),
			Token:        accessToken,
			RefreshToken: refreshToken,
		}
		_ = loginSuccess.WriteToResponse(w, response)

		return
	}
}
