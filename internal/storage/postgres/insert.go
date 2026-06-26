package postgres

import (
	"LogSentry/internal/models"
	"database/sql"
)

func insertEntry(db *sql.DB, log models.LogEntry) error {
	query := `
	INSERT INTO log_entries
	(timestamp, category, source, details)
	VALUES ($1, $2, $3, $4)
	`

	_, err := db.Exec(
		query,
		log.TimeStamp,
		log.Category,
		log.Source,
		log.Details,
	)

	return err
}

func InsertLogs(db *sql.DB, report models.LogReport) error {

	for _, log := range report.Errors {
		if err := insertEntry(db, log); err != nil {
			return err
		}
	}

	for _, log := range report.Warns {
		if err := insertEntry(db, log); err != nil {
			return err
		}
	}

	for _, log := range report.Infos {
		if err := insertEntry(db, log); err != nil {
			return err
		}
	}

	for _, log := range report.Unknown {
		if err := insertEntry(db, log); err != nil {
			return err
		}
	}

	return nil
}