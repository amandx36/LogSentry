	package routes 

	import (
		"database/sql"
		"LogSentry/internal/api/controller"
		"github.com/gin-gonic/gin"
	)
	func SearchByKeyWordsEP(rte *gin.Engine, db *sql.DB) {
		rte.GET("/search/keywords", controller.SearchByKeyWord(db))
	}