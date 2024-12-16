package utils

import (
	"encoding/json"
	"net/http"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/logger"
)

type ReqTime string

const ReqStartTime ReqTime = "reqStartTime"

func RespondWithError(w http.ResponseWriter, code int, msg string) {

	if code > 499 {
		logger.Error("Responding with 5XX error: %s", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}
