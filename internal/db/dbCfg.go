package databaseCfg

import (
	"database/sql"
	"fmt"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/db/database"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	DB   *database.Queries
	Conn *sql.DB
}

func NewDBConfig(dbURL string) (*DBConfig, error) {

	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL environment variable is not set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("can't connect to the database: %v", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping database: %v", err)
	}

	db := database.New(conn)

	dbConfig := &DBConfig{
		DB:   db,
		Conn: conn,
	}

	return dbConfig, nil
}
