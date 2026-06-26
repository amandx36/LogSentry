package postgres

import (
	"database/sql"
	"fmt"
	"LogSentry/internal/config"
	// Add this blank import right here!
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect( cfg config.Config ) (*sql.DB, error) {
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

	
	err = createTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return db, nil 
}


func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS log_entries (
		id SERIAL PRIMARY KEY,
		timestamp TEXT NOT NULL,
		category VARCHAR(10) NOT NULL,
		source VARCHAR(100),
		details TEXT NOT NULL
	);
	`
	
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	fmt.Println("Table created ! ")
	return nil
	
}