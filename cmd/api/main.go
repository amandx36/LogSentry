package main
import (
	"github.com/gin-gonic/gin"
	"LogSentry/internal/api/routes"	
)
func main() {
	api := gin.Default()

routes.HealthCheck(api)
routes.ApiVersion(api)

api.Run(":8080")
}