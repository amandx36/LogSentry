package controller

import (
	"LogSentry/internal/services/analytics"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func GetAnalyticsEP(db *sql.DB) gin.HandlerFunc {
	return func (ctx *gin.Context) {
		analytics , err := analytics.GetAnalytics(db)
		if err != nil{
			ctx.JSON(500, gin.H{
				"error": "failed to get analytics",
			})
			return
		}
		ctx.JSON(200, analytics)
	}
}	