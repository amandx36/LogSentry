package routes 

import (
	"database/sql"
	"LogSentry/internal/api/controller"
	"github.com/gin-gonic/gin"	
)
// curl "http://localhost:8080/search/date?start_date=2026-06-24&end_date=2026-06-26"
func SearchByDateEP(rte *gin.Engine, db *sql.DB) {
	rte.GET("/search/date", controller.SearchByDateEP(db))
}