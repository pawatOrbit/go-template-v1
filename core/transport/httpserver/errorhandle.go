package httpserver

import (
	"encoding/json"
	"net/http"
)

type ModelResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HandleInternalServerError(w http.ResponseWriter, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	resp := ModelResp{
		Status:  httpStatusCode,
		Message: "Internal Server Error",
	}

	json.NewEncoder(w).Encode(resp)
}
