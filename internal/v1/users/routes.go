package users

import (
	"net/http"
)

// NewRouter returns a new http.ServeMux with v1 routes configured
func NewRouter() *http.ServeMux {
	userRouter := http.NewServeMux()

	userRouter.HandleFunc("POST /sendMail", MailHandler)
	userRouter.HandleFunc("POST /sendHTML", HtmlMailHandler)

	return userRouter
}
