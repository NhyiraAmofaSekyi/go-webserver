package v1

import (
	"net/http"

	utils "github.com/NhyiraAmofaSekyi/go-webserver/utils"
)

// Handler function for the "healthz" endpoint
func healthzHandler(w http.ResponseWriter, r *http.Request) {

	utils.RespondWithJSON(w, 200, map[string]string{"status": "ok", "route": "v1"})
}

func secureHandler(w http.ResponseWriter, r *http.Request) {

	utils.RespondWithJSON(w, 200, map[string]string{"status": "ok", "route": "secure"})
}
