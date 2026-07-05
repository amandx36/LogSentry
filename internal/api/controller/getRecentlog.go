package controller

import (
	"database/sql"
	"net/http"
	"strconv"

	"LogSentry/internal/services/search"
	"github.com/gin-gonic/gin"
)

func GetRecentLog(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		limitStr := ctx.Query("limit")

		if limitStr == "" {
			limitStr = "10" 
		}
		// converting ascii ot int 
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "limit must be an integer",
			})
			return
		}

		logs, err := search.GetRecentLogs(db, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, logs)
	}
}