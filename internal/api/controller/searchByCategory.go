package controller

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"LogSentry/internal/services/search"
)

func SearchByCategory(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context){
		category := ctx.Query("category")
		if category == "" {
			ctx.JSON(400,
			gin.H{
				"error":"category query parameter is required",
			},
			)
			return 
		}
		logs , err := search.SearchByCategory(db, category)
		 if err != nil {
			ctx.JSON(500,
			gin.H{
				"Error":err.Error(),
			})
			return 
		 }
		 ctx.JSON(200, logs )
		 
	}
}