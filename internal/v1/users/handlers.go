package users

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

func FileForm(w http.ResponseWriter, r *http.Request) {
	// Define the endpoint where the form will submit the data

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	api := "/api/v1/"
	endpoint := "http://" + host + ":" + port + api + "users/upload"

	// Parse the template file
	tmpl, err := template.ParseFiles("./internal/templates/file_form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Data to pass to the template
	data := struct {
		Endpoint string
	}{
		Endpoint: endpoint,
	}

	// Execute the template with the data
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func Upload(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20) // 10MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, 500, "Error retrieving the file")
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.RespondWithError(w, 500, "Error reading the file")
		return
	}

	fileSize := len(fileBytes)

	fileType := strings.ToLower(filepath.Ext(handler.Filename))

	response := map[string]interface{}{
		"fileName": handler.Filename,
		"fileType": fileType,
		"fileSize": fileSize,
	}

	utils.RespondWithJSON(w, 200, response)
}
