package models

import (
	"database/sql"
	"fmt"
)

func CreateTables(db *sql.DB) error {
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
