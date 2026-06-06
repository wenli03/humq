package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
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
	database.SeedDefaultAdmin(db)

	r := api.SetupRouter(db, cfg)

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "web/dist"
	}
	if info, err := os.Stat(staticDir); err == nil && info.IsDir() {
		fs := http.FileServer(http.Dir(staticDir))
		r.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"code": 4004, "msg": "not found"})
				return
			}
			fs.ServeHTTP(c.Writer, c.Request)
		})
		log.Printf("Serving static files from %s", staticDir)
	}

	log.Printf("HU MQ server starting on :%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
