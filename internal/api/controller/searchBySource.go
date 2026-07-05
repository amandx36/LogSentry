package controller

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"LogSentry/internal/services/search"
)

func SearchBySource(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context){
		source := ctx.Query("source")
		if source == "" {
			ctx.JSON(400,
			gin.H{
				"error":"source query parameter is required",
			},
			)
			return 
		}
		logs , err := search.SearchBySource(db, source)
		 if err != nil {
			ctx.JSON(500,
			gin.H{
				"Error":err.Error(),
			})
			return 
		 }
		ctx.JSON(http.StatusOK, gin.H{
    "count": len(logs),
    "data":  logs,
})
		 
	}
}