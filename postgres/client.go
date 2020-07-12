package postgres

import (
	"database/sql"
	"fmt"
	"os"
)

// NewConnection return a new connection for the Postgres database
func NewConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_SSLMODE"))

	db, err := sql.Open(os.Getenv("DATABASE_DRIVER"), connStr)

	if err != nil {
		return &sql.DB{}, err
	}

	if err = db.Ping(); err != nil {
		return &sql.DB{}, err
	}

	return db, nil
}
