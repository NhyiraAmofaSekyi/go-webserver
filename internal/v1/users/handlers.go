package users

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/config"
	utils "github.com/NhyiraAmofaSekyi/go-webserver/utils"
	aws "github.com/NhyiraAmofaSekyi/go-webserver/utils/aws/awsS3"
	email "github.com/NhyiraAmofaSekyi/go-webserver/utils/email"
	uuid "github.com/google/uuid"
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
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"message": "server error"})
		return
	}

	err = email.SendMail(params.Subject, params.Email, params.Name)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"message": "failed to send email"})
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
		utils.RespondWithError(w, http.StatusInternalServerError, "error parsing json")
		return
	}

	err = email.SendHTML(params.Subject, params.Email, params.Name)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error sending email")
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

func ListObj(w http.ResponseWriter, r *http.Request) {

	err := aws.ListBucketOBJ()

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func GetObj(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Key string `json:"key"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"message": "server error"})
		return
	}

	_, err = aws.GetObject(params.Key, "arn:aws:s3:eu-north-1:049991758581:accesspoint/test2")
	bucket := os.Getenv("AWS_BUCKET")
	region := os.Getenv("AWS_BUCKET_REGION")
	url := "https://" + bucket + ".s3." + region + ".amazonaws.com/" + params.Key
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, url)
}

func Upload(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20) // 10MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving the file")
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error reading the file")
		return
	}
	openFile := handler.Filename

	fileSize := len(fileBytes)

	fileType := strings.ToLower(filepath.Ext(handler.Filename))
	println("file type", fileType)

	id := uuid.New()

	bucket := os.Getenv("AWS_BUCKET")
	region := os.Getenv("AWS_BUCKET_REGION")

	contentType := mime.TypeByExtension(fileType)
	println("content type", contentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	key := id.String() + fileType

	err = aws.UploadFile(bucket, key, openFile, contentType)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error uploading")
		return
	}

	url := "https://" + bucket + ".s3." + region + ".amazonaws.com/" + key
	response := map[string]interface{}{
		"fileName": handler.Filename,
		"fileType": fileType,
		"fileSize": fileSize,
		"url":      url,
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	dbConfig := config.AppConfig.DBConfig

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "bad request"})
		return
	}

	user, err := dbConfig.DB.CreateUser(r.Context(), params.Name)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)

}
func GetUsers(w http.ResponseWriter, r *http.Request) {

	dbConfig := config.AppConfig.DBConfig

	users, err := dbConfig.DB.GetUsers(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, users)

}
