package config

import (
	"log"
	"sync"

	databaseCfg "github.com/NhyiraAmofaSekyi/go-webserver/internal/db"
)

// "golang.org/x/oauth2"

var (
	AppConfig *Config
	once      sync.Once
)

type Config struct {
	DBConfig *databaseCfg.DBConfig
}

func Initialize(dbConfig *databaseCfg.DBConfig) {
	once.Do(func() {
		if dbConfig == nil {
			log.Fatal("Database configuration is nil")
		}
		AppConfig = &Config{
			DBConfig: dbConfig,
		}
	})
}

// func NewAuth2() oauth2.Config {

// 	if err := godotenv.Load(); err != nil {
// 		log.Println("No .env file found, assuming production settings")
// 	}
// 	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
// 	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

// 	log.Println("loaded auth variables")
// 	AppConfig.GoogleLoginConfig = oauth2.Config{
// 		RedirectURL:  "http://localhost:8080/api/v1/auth/google/callback",
// 		ClientID:     googleClientID,
// 		ClientSecret: googleClientSecret,
// 		Scopes: []string{
// 			"https://www.googleapis.com/auth/userinfo.email",
// 			"https://www.googleapis.com/auth/userinfo.profile",
// 		},
// 		Endpoint: google.Endpoint,
// 	}

// 	return AppConfig.GoogleLoginConfig

// }
