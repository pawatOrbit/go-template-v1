package middleware

import (
	"encoding/json"
	"net/http"
)

type ModelResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	resp := ModelResp{
		Status:  http.StatusNotFound,
		Message: "Not Found",
	}

	json.NewEncoder(w).Encode(resp)
}
