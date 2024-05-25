package databaseCfg

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/db/database"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	DB *database.Queries
}

func NewDBConfig() (*DBConfig, error) {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DB_URL not found in environment")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("can't connect to the database: %v", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping database: %v", err)
	}

	db := database.New(conn)

	dbConfig := &DBConfig{
		DB: db,
	}
	fmt.Printf("DBCfg.DB: %v\n", dbConfig.DB)

	return dbConfig, nil
}
