package routes 
import (
	"LogSentry/internal/api/controller"
	"github.com/gin-gonic/gin"
	"database/sql"
)

// eg curl "http://localhost:8080/search/category?category=ERROR"
// data is send thorugh response body in json format 
func SearchByCategoryEP(rte *gin.Engine, db *sql.DB) {
	rte.GET("/search/category", controller.SearchByCategory(db))
}

