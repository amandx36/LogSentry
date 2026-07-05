package routes 
import (
	"LogSentry/internal/api/controller"
	"github.com/gin-gonic/gin"
	"database/sql"
)

func SearchByCategoryEP(rte *gin.Engine, db *sql.DB) {
	rte.GET("/search/:category", controller.SearchByCategory(db))
}