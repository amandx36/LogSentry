package search

import (
	"LogSentry/internal/models"
	"database/sql"
	
)

func SearchByCategory(db *sql.DB, category string) ([]models.LogEntry, error) {
	query := `SELECT id, timestamp, category, source, details FROM log_entries WHERE category = $1;`
	return fetchLogs(db, query, category)
}

func SearchBySource(db *sql.DB, source string) ([]models.LogEntry, error) {
	query := `SELECT id, timestamp, category, source, details FROM log_entries WHERE source = $1;`
	return fetchLogs(db, query, source)
}

func SearchByKeywords(db *sql.DB, keyword string) ([]models.LogEntry, error) {
	query := `SELECT id, timestamp, category, source, details FROM log_entries WHERE details ILIKE '%' || $1 || '%';`
	return fetchLogs(db, query, keyword)
}

func GetRecentLogs(db *sql.DB, limit int) ([]models.LogEntry, error) {
	query := `SELECT id, timestamp, category, source, details FROM log_entries ORDER BY id DESC LIMIT $1;`
	return fetchLogs(db, query, limit)
}

func SearchByDate(db *sql.DB, startDate string, endDate string) ([]models.LogEntry, error) {
	query := `SELECT id, timestamp, category, source, details FROM log_entries WHERE timestamp >= $1 AND timestamp <= $2;`
	return fetchLogs(db, query, startDate, endDate)
}