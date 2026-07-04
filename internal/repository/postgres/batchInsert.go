package postgres

import (
	"LogSentry/internal/models"
	"database/sql"
	"fmt"
	"strings"
)

func BatchInsert(db *sql.DB, logs []models.LogEntry) error {

	if len(logs) == 0 {
		return nil
	}

	// Start Transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Rollback if anything fails
	defer tx.Rollback()

	var (
		values      []string
		args        []interface{}
		placeholder = 1
	)

	for _, log := range logs {

		values = append(values,
			fmt.Sprintf("($%d,$%d,$%d,$%d)",
				placeholder,
				placeholder+1,
				placeholder+2,
				placeholder+3,
			),
		)

		args = append(args,
			log.TimeStamp,
			log.Category,
			log.Source,
			log.Details,
		)

		placeholder += 4
	}

	query := fmt.Sprintf(`
	INSERT INTO log_entries
	(timestamp, category, source, details)
	VALUES %s
	`, strings.Join(values, ","))

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	// Commit Transaction
	return tx.Commit()
}
