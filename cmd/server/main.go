package main

import (
	"log"

	"LogSentry/internal/api/routes"
	"LogSentry/internal/config"
	"LogSentry/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Loadconfig("internal/config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	api := gin.Default()

	routes.HealthCheckEP(api)
	routes.ApiVersionEP(api)
	routes.DashBoardEP(api, db)
	routes.SearchByCategoryEP(api,db)
	routes.SearchBySourceEP(api,db)
	routes.SearchByKeyWordsEP(api,db)
	routes.GetRecentLogEP(api,db)
	routes.SearchByDateEP(api, db)

	if err := api.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}