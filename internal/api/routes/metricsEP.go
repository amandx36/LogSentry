package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GetMatrixDetails(rte *gin.Engine) {
	rte.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
