package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegestAllRoutes(rte *gin.Engine, db *sql.DB) {
	SearchByDateEP(rte, db)
	SearchBySourceEP(rte, db)
	SearchByKeyWordsEP(rte, db)
	GetRecentLogEP(rte, db)
	AnalyticsEP(rte, db)
	DashBoardEP(rte, db)
	HealthCheckEP(rte)
	ApiVersionEP(rte)
	SearchByCategoryEP(rte, db)
	GetMatrixDetails(rte)
}
