package routes

import (
    "database/sql"

    "LogSentry/internal/api/controller"

    "github.com/gin-gonic/gin"
)

func DashBoardEP(r *gin.Engine, db *sql.DB) {
    r.GET("/dashboard", controller.GetDashboard(db))
}