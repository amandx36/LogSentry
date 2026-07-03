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

	if err := BatchInsert(db, report.Errors); err != nil {
		return err
	}

	if err := BatchInsert(db, report.Warns); err != nil {
		return err
	}

	if err := BatchInsert(db, report.Infos); err != nil {
		return err
	}

	if err := BatchInsert(db, report.Unknown); err != nil {
		return err
	}

	return nil
}