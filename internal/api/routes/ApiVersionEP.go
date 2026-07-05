package routes

import (
	"LogSentry/internal/api/controller"

	"github.com/gin-gonic/gin"
)

func ApiVersionEP (route *gin.Engine){
	route.GET("api/version",controller.Apiversion)
}