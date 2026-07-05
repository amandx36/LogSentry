package routes


import (
	"LogSentry/internal/api/controller"

	"github.com/gin-gonic/gin"
)

func HealthCheckEP(r *gin.Engine) {
	r.GET("/ping", controller.Ping)
}