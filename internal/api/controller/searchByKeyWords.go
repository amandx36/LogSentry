package controller

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"LogSentry/internal/services/search"
)

func SearchByKeyWord(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context){
		keyword := ctx.Query("keyword")
		if keyword == "" {
			ctx.JSON(400,
			gin.H{
				"error":"keyword query parameter is required",
			},
			)
			return 
		}
		logs , err := search.SearchByKeywords(db, keyword)
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