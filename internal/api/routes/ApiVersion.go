package routes

import (
	"LogSentry/internal/api/controller"

	"github.com/gin-gonic/gin"
)

func ApiVersion (route *gin.Engine){
	route.GET("api/version",controller.Apiversion)
}