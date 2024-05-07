package users

import (
	"fmt"
	"net/http"

	utils "github.com/NhyiraAmofaSekyi/go-webserver/utils/email"
)

func MailHandler(w http.ResponseWriter, r *http.Request) {
	err := utils.SendMail()
	if err != nil {
		http.Error(w, "Failed to send mail", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Mail sent successfully")
}

func HtmlMailHandler(w http.ResponseWriter, r *http.Request) {
	err := utils.SendHTML("Welcome!")
	if err != nil {
		http.Error(w, "Failed to send HTML mail", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "HTML mail sent successfully")
}
