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
		Email   string `json:"email"`
		Subject string `json:"subject"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithJSON(w, 500, map[string]string{"message": "server error"})
		return
	}

	err = email.SendMail(params.Subject, params.Email, params.Name)
	if err != nil {
		utils.RespondWithJSON(w, 400, map[string]string{"message": "failed to send email"})
		return
	}
	fmt.Fprintln(w, "Mail sent successfully")
}

func HtmlMailHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Subject string `json:"subject"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithJSON(w, 500, map[string]string{"message": "server error"})
		return
	}

	err = email.SendHTML(params.Subject, params.Email, params.Name)
	if err != nil {
		utils.RespondWithJSON(w, 400, map[string]string{"message": "failed to send email"})
		return
	}
	fmt.Fprintln(w, "HTML mail sent successfully")
}
