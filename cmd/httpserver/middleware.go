package httpserver

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"orderservice/state"
	utils "orderservice/utils"
	"strings"
)

func AuthMiddleware(app *state.State) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractTokenFromHeader(r)
			if tokenStr == "" {
				app.Logger.PrintError(fmt.Errorf("no token provided"), map[string]string{
					"context": "authorization",
				})
				_ = Unauthorized.WriteToResponse(w, nil)
				return
			}

			var claims utils.Claims
			token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.Config.SecretKey), nil
			})

			if err != nil || !token.Valid {
				app.Logger.PrintError(fmt.Errorf("invalid token"), map[string]string{
					"context": "authorization",
				})
				_ = Unauthorized.WriteToResponse(w, nil)
				return
			}
			ctx := context.WithValue(r.Context(), "userid", claims.UserID.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userid").(string)
	return userID, ok
}
