package httpserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"net/http"
	"orderservice/repository"
	"orderservice/state"
	utils "orderservice/utils"
)

type RegistrationRequestPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func HandleRegisterUser(app *state.State) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		request := RegistrationRequestPayload{}
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
			if errors.Is(err, sql.ErrNoRows) {
				app.Logger.PrintInfo(fmt.Sprintf("user with email %s not found", request.Email), map[string]string{})
			} else {
				app.Logger.PrintError(err, map[string]string{
					"context": "Error fetching user by email",
				})
				_ = InternalError.WriteToResponse(w, nil)
				return
			}
		}

		if user != nil {
			app.Logger.PrintInfo(fmt.Sprintf("User already exists: %s", request.Email), map[string]string{
				"context": "user registration",
			})
			_ = UserAlreadyExist.WriteToResponse(w, nil)
			return
		}

		passwordHash, err := utils.HashPassword(request.Password)

		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Failed to hash password",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}
		userID, err := uuid.NewV4()

		if err != nil {
			_ = InternalError.WriteToResponse(w, nil)
			return
		}

		user = &repository.User{
			ID:       userID,
			Name:     request.Name,
			Email:    request.Email,
			Password: passwordHash,
		}

		if err = app.Repository.CreateUser(ctx, user); err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Failed to create user",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}

		_ = UserCreated.WriteToResponse(w, nil)
		return

	}
}
