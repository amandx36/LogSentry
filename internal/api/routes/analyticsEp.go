package routes 
import (
	"database/sql"
	"LogSentry/internal/api/controller"
	"github.com/gin-gonic/gin"	
)
// curl "http://localhost:8080/analytics"
func AnalyticsEP(rte *gin.Engine, db *sql.DB) {
	rte.GET("/analytics", controller.GetAnalyticsEP(db))	
}