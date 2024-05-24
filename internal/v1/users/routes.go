package users

import (
	"net/http"

	databaseCfg "github.com/NhyiraAmofaSekyi/go-webserver/internal/db"
)

func NewRouter(dbCFG *databaseCfg.DBConfig) *http.ServeMux {
	userRouter := http.NewServeMux()

	userRouter.HandleFunc("POST /sendMail", MailHandler)
	userRouter.HandleFunc("POST /sendHTML", HtmlMailHandler)
	userRouter.HandleFunc("/fileForm", FileForm)
	userRouter.HandleFunc("/upload", Upload)
	userRouter.HandleFunc("/listObj", ListObj)
	userRouter.HandleFunc("/getObj", GetObj)
	userRouter.HandleFunc("/createUser", CreateUser(dbCFG))

	return userRouter
}
