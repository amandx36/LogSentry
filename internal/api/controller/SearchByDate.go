package controller

import (
	"database/sql"
	"LogSentry/internal/services/search"
	"github.com/gin-gonic/gin"
)

func SearchByDateEP(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startDate := ctx.Query("start_date")
		endDate := ctx.Query("end_date")

		if startDate == "" || endDate == "" {
			ctx.JSON(400, gin.H{
				"error": "start_date and end_date query parameters are required",
			})
			return
		}

		logs, err := search.SearchByDate(db, startDate, endDate)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": "failed to search logs by date",
			})
			return
		}

		ctx.JSON(200, logs)
	}
}