package search
import(
	"LogSentry/internal/models"
	"database/sql"
	"fmt"
)

func fetchLogs(db *sql.DB, query string, args ...any) ([]models.LogEntry, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var logs []models.LogEntry

	for rows.Next() {
		var log models.LogEntry
		var id int 
		
		err := rows.Scan(&id, &log.TimeStamp, &log.Category, &log.Source, &log.Details)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return logs, nil
}

