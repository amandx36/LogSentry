package analytics

import (
	"database/sql"
	"fmt"
)

func GetTotalLogs(db *sql.DB) int {
    query := `SELECT COUNT(*) FROM log_entries`

    var totalLogs int

    err := db.QueryRow(query).Scan(&totalLogs)
    if err != nil {
        fmt.Println("Error while fetching total logs:", err)
        return 0
    }

    return totalLogs
}

// for specific categories count 
func GetSpecificCategoriesCount(db *sql.DB, cat string) int {
    query := `SELECT COUNT(*) FROM log_entries WHERE category=$1;`
    var specificCategoriesCount int
    
    err := db.QueryRow(query, cat).Scan(&specificCategoriesCount)
    if err != nil {
        fmt.Printf("Error while getting specific category count: %v\n", err)
        return 0 
    }

    return specificCategoriesCount
}

// geting top source 
func GetTopSource(db *sql.DB) (string, int) {
query := `SELECT source, COUNT(*) FROM log_entries GROUP BY source ORDER BY COUNT(*) DESC LIMIT 1`

	var source string
	var count int
	err := db.QueryRow(query).Scan(&source, &count)
	if err != nil {
		fmt.Println("Error fetching top source:", err)
		return "", 0
	}

	return source, count
}
// getting most frequent error 
func GetMostFrequentError(db *sql.DB) (string, int) {
query := `SELECT details, COUNT(*) FROM log_entries WHERE category='ERROR' GROUP BY details ORDER BY COUNT(*) DESC LIMIT 1`
	var message string
	var count int
	err := db.QueryRow(query).Scan(&message, &count)
	if err != nil {
		fmt.Println("Error fetching most frequent error:", err)
		return "", 0
	}

	return message, count
}

func GetErrorRate(db *sql.DB) float64 {
	totalLogs := GetTotalLogs(db)
	if totalLogs == 0 {
		return 0.0
	}

	errorLogs := GetSpecificCategoriesCount(db, "ERROR")
	
	errorRate := (float64(errorLogs) / float64(totalLogs)) * 100.0
	return errorRate
}