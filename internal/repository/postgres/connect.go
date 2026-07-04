package postgres

import (
	"LogSentry/internal/config"
	"LogSentry/internal/models"
	"database/sql"
	"fmt"
	// Add this blank import right here!
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(cfg config.Config) (*sql.DB, error) {
	// data base url
	ConnString := cfg.DatabaseUrl

	db, err := sql.Open("pgx", ConnString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// ping to check connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	fmt.Println("Connected successfully to the PostgreSQL database!")

	err = models.CreateTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return db, nil
}
