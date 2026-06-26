package dashboard

import (
	"LogSentry/models"
	"fmt"
)

func DashViewer(report models.DashBoardDetails) {
	fmt.Println("Total Logs :", report.TotalLogs)
	fmt.Println("Errors     :", report.Errors)
	fmt.Println("Warns      :", report.Warns)
	fmt.Println("Infos      :", report.Infos)
	fmt.Println("Unknown    :", report.Unknown)
}