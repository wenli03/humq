package main

import (
	"log"

	"github.com/wenli03/humq/internal/api"
	"github.com/wenli03/humq/internal/config"
	"github.com/wenli03/humq/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DB.DSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	database.Migrate(db)

	r := api.SetupRouter(db, cfg)
	log.Printf("HU MQ server starting on :%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
