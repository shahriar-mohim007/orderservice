package httpserver

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"orderservice/state"
)

func HandleCancelOrder(app *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		consignmentId := chi.URLParam(req, "id")
		err := app.Repository.CancelOrder(ctx, consignmentId)
		if err != nil {
			_ = OrderCancelError.WriteToResponse(w, nil)
			return
		}
		_ = OrderCancelSuccess.WriteToResponse(w, nil)
		return
	}
}
