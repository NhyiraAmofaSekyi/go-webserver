package v1

import (
	"log"
	"net/http"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/middleware"
	utils "github.com/NhyiraAmofaSekyi/go-webserver/utils"
)

// Handler function for the "healthz" endpoint
func HealthzHandler(w http.ResponseWriter, r *http.Request) {

	utils.RespondWithJSON(w, 200, map[string]string{"status": "ok", "route": "v1"})
}

func SecureHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.AuthUserID).(string)
	log.Println("user logged in: ", userID)
	utils.RespondWithJSON(w, 200, map[string]string{"status": "ok", "route": "secure", "userID": userID})
}
