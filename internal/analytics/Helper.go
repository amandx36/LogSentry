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
func GetTotalCategories(db *sql.DB)int{
	query :=`SELECT category, COUNT(*) FROM log_entries GROUP BY category;`
    var totalCategories int 
    err := db.QueryRow(query).Scan(&totalCategories)
    if err != nil{
        fmt.Println("Error while Fetching the Categories %w",err)
    
    }
    return totalCategories;
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
