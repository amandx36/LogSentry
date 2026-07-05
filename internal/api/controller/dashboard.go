package controller

import (
	"LogSentry/internal/services/dashboard"
	"database/sql"
	"github.com/gin-gonic/gin"
)

// StatsInfo, err := dashboard.GetOverallStats(db)
// defer db.Close()
// if err != nil {
// 	fmt.Println("Error while getting the info from the database :)")
// 	return
// }
// fmt.Println("Total Errors ", StatsInfo.Errors)
// fmt.Println("Total Info ", StatsInfo.Infos)
// fmt.Println("Total Total Logs  ", StatsInfo.TotalLogs)
// fmt.Println("Total Unknown  ", StatsInfo.Unknown)
// fmt.Println("Total Warns  ", StatsInfo.Warns)

// problem don't need the global db connection so i use the closer and return the handler function to router

func GetDashboard(db *sql.DB) gin.HandlerFunc {

    return func(ctx *gin.Context) {

        stats, err := dashboard.GetOverallStats(db)

        if err != nil {
            ctx.JSON(500, gin.H{
                "error": err.Error(),
            })
            return
        }

        ctx.JSON(200, stats)
    }
}