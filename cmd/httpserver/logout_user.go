package httpserver

import (
	"net/http"
	"orderservice/state"
	"sync"
	"time"
)

var blacklistedTokens sync.Map

func HandleLogout(app *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		token := ExtractTokenFromHeader(req)
		if token == "" {
			_ = Unauthorized.WriteToResponse(w, nil)
			return
		}

		expiration := time.Now().Add(24 * time.Hour)
		blacklistedTokens.Store(token, expiration)

		_ = LogoutSuccess.WriteToResponse(w, nil)
		return
	}
}
