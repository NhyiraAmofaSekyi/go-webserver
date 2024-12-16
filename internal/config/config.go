package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"

	databaseCfg "github.com/NhyiraAmofaSekyi/go-webserver/internal/db"
	"github.com/NhyiraAmofaSekyi/go-webserver/internal/logger"
)

// "golang.org/x/oauth2"

var (
	Config *Con
	once   sync.Once
	env    string
)

type ServerConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	ClientURL string `yaml:"client_url"`
	Debug     bool   `yaml:"debug"`
}
type YAMLConfig struct {
	Environments struct {
		Local struct {
			Server ServerConfig `yaml:"server"`
		} `yaml:"local"`
		Production struct {
			Server ServerConfig `yaml:"server"`
		} `yaml:"production"`
	} `yaml:"environments"`
}

type Con struct {
	DBConfig  *databaseCfg.DBConfig
	ClientURL string `yaml:"client_url"`
	APIHost   string `yaml:"api_host"`
	APIPort   int    `yaml:"api_port"`
}

func Initialise() {
	once.Do(func() {
		envPtr := flag.String("env", "development", "Define the application environment (development or production)")
		flag.Parse()
		env = *envPtr

		data, err := os.ReadFile("internal/config/config.yaml")
		if err != nil {
			log.Fatal("Error reading config file:", err)
		}
		envFile := fmt.Sprintf(".env.%s", env)
		if err := godotenv.Load(envFile); err != nil {
			log.Fatalf("error loading env file %s: %v", envFile, err)
			return
		}

		var yamlConfig YAMLConfig
		if err := yaml.Unmarshal(data, &yamlConfig); err != nil {
			log.Fatal("Error parsing config:", err)
		}
		var serverConfig ServerConfig
		if env == "production" {
			serverConfig = yamlConfig.Environments.Production.Server
		} else {
			serverConfig = yamlConfig.Environments.Local.Server
		}

		// Initialize logger with debug mode from config
		logger.Init(serverConfig.Debug)
		logger.Debug("Initializing configuration for environment: %s", env)

		dbURL := os.Getenv("DB_URL")
		if dbURL == "" {
			logger.Fatal("DB_URL environment variable is not set in %s", envFile)
		}
		logger.Debug("Database URL loaded from environment")

		dbConfig, err := databaseCfg.NewDBConfig(dbURL)
		if err != nil {
			logger.Fatal("Database initialization failed: %v", err)
		}
		if dbConfig == nil {
			logger.Fatal("Database configuration is nil")
		}
		logger.Info("Database connection established successfully")

		Config = &Con{
			DBConfig:  dbConfig,
			ClientURL: serverConfig.ClientURL,
			APIHost:   serverConfig.Host,
			APIPort:   serverConfig.Port,
		}
		logger.Debug("Base configuration created: %+v", Config)

		// Override with environment variables if present
		if url := os.Getenv("CLIENT_URL"); url != "" {
			Config.ClientURL = url
			logger.Debug("Overrode ClientURL from environment: %s", url)
		}
		if host := os.Getenv("API_HOST"); host != "" {
			Config.APIHost = host
			logger.Debug("Overrode APIHost from environment: %s", host)
		}
	})

}
