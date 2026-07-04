package dashboard

import (
	"LogSentry/internal/models"
	"database/sql"
	"fmt"
)

func DashViewer(report models.DashBoardDetails) {
	fmt.Println("Current Uploaded dashboard :)")

	fmt.Println("Total Logs :", report.TotalLogs)
	fmt.Println("Errors     :", report.Errors)
	fmt.Println("Warns      :", report.Warns)
	fmt.Println("Infos      :", report.Infos)
	fmt.Println("Unknown    :", report.Unknown)
}
func GetOverallStats(db *sql.DB) (models.DashBoardDetails, error) {
	var stats models.DashBoardDetails
	query := `SELECT category, COUNT(*) FROM log_entries GROUP BY category;`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error in executing the error %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var category string
		var count int
		// Scan takes the columns from the current database row and puts them into our Go variables
		if err := rows.Scan(&category, &count); err != nil {
			return stats, fmt.Errorf("failed to read row data: %w", err)
		}

		switch category {
		case "ERROR":
			stats.Errors = count
		case "WARN":
			stats.Warns = count
		case "INFO":
			stats.Infos = count
		default:
			stats.Unknown += count
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error while executing the loop on the rows %w", err)
	}

	stats.TotalLogs = stats.Errors + stats.Warns + stats.Infos + stats.Unknown
	return stats, nil

}
