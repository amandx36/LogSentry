package main

import (
	"log"
	"time"

	"LogSentry/internal/api/routes"
	"LogSentry/internal/config"
	"LogSentry/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	start := time.Now()



	step := time.Now()
	cfg, err := config.Loadconfig("internal/config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[Startup] Config loaded in %v", time.Since(step))

	step = time.Now()
	db, err := postgres.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Printf("[Startup] Database connected in %v", time.Since(step))

	step = time.Now()
	api := gin.Default()
	log.Printf("[Startup] Gin engine created in %v", time.Since(step))

	step = time.Now()
	routes.RegestAllRoutes(api, db)
	log.Printf("[Startup] Routes registered in %v", time.Since(step))

	log.Printf("[Startup] Total startup time: %v", time.Since(start))
	log.Println("[Startup] Server listening on http://localhost:8080")

	if err := api.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}