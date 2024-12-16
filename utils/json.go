package utils

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

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

func RespondWithJSONContext(ctx context.Context, w http.ResponseWriter, code int, payload interface{}) {
	// Retrieve the start time from the context
	startTime, ok := ctx.Value(ReqStartTime).(time.Time)
	if !ok {
		log.Println("Could not retrieve request start time from context")
	}

	// Marshal the payload into JSON
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ok {
		duration := time.Since(startTime)
		log.Printf("Request took %v", duration)
	}

	// Set the headers and write the response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}
