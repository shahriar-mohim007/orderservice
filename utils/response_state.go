package utilis

import (
	"encoding/json"
	"net/http"
)

type CommonResponse struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseState struct {
	StatusCode int
	Message    string
	Type       string
}

func (rs ResponseState) WriteToResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs.StatusCode)
	return json.NewEncoder(w).Encode(CommonResponse{
		Message: rs.Message,
		Type:    rs.Type,
		Code:    rs.StatusCode,
		Data:    data,
	})
}
