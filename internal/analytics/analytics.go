package analytics

import (
	"LogSentry/internal/models"
	"database/sql"
	"time"
)
func 	GetAnalytics(db *sql.DB)(models.Analytics,error){
	analytics := models.Analytics{}

analytics.TotalLogs = GetTotalLogs(db)

analytics.TotalErrors = GetSpecificCategoriesCount(db, "ERROR")

analytics.TotalWarns = GetSpecificCategoriesCount(db, "WARN")

analytics.TotalInfos = GetSpecificCategoriesCount(db, "INFO")

analytics.TotalUnknown = GetSpecificCategoriesCount(db, "UNKNOWN")

analytics.ErrorRate = GetErrorRate(db)

analytics.TopSource, analytics.TopSourceCount = GetTopSource(db)

analytics.GeneratedAt = time.Now()

return analytics, nil
}
