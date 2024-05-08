package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	utils "github.com/NhyiraAmofaSekyi/go-webserver/utils"
	email "github.com/NhyiraAmofaSekyi/go-webserver/utils/email"
)

func MailHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name    string `json:"name"`
		Subject string `json:"subject"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithJSON(w, 400, fmt.Sprintf("Error passing json: %v", err))
		return
	}

	err = email.SendMail(params.Subject, params.Name)
	if err != nil {
		http.Error(w, "Failed to send mail", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Mail sent successfully")
}

func HtmlMailHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name    string `json:"name"`
		Subject string `json:"subject"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithJSON(w, 400, fmt.Sprintf("Error passing json: %v", err))
		return
	}

	err = email.SendHTML(params.Subject, params.Name)
	if err != nil {
		utils.RespondWithJSON(w, 400, fmt.Sprintf("Failed to send HTML mail %v", err))
		return
	}
	fmt.Fprintln(w, "HTML mail sent successfully")
}
