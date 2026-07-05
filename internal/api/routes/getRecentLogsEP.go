package routes 
import (
	"LogSentry/internal/api/controller"
	"database/sql"
	"github.com/gin-gonic/gin"	
)
func GetRecentLogEP(rte *gin.Engine, db *sql.DB) {
	rte.GET("/logs/recent", controller.GetRecentLog(db))
}