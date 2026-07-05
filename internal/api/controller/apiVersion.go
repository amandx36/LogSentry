package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Apiversion(ctg *gin.Context) {
	ctg.JSON(
		http.StatusOK,
		gin.H{
			"Version": "1.0.0",
			"EndPoints":"Working",	
		},
	)
}