package main

import (
	"log"
	"net/http"

	v1 "github.com/NhyiraAmofaSekyi/go-webserver/internal/v1"
	"github.com/NhyiraAmofaSekyi/go-webserver/utils"
)

func main() {

	mainMux := http.NewServeMux()

	log.Println("server running on port 8080")
	mainMux.Handle("/v1/", http.StripPrefix("/v1", v1.NewRouter()))

	corsEnabledMux := utils.CorsWrapper(mainMux)

	srv := &http.Server{
		Handler: corsEnabledMux, // Your wrapped handler
		Addr:    ":8080",        // Listen address
		// Other configurations like ReadTimeout, WriteTimeout, etc.
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
