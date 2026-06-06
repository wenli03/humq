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

	db, err := database.Connect(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		log.Println("WARN: database unavailable, running in standalone mode")
		db = nil
	}

	if db != nil {
		database.Migrate(db)
		database.SeedDefaultAdmin(db)
		log.Println("Database connected")
	} else {
		log.Println("Standalone mode - no database required")
	}

	r := api.SetupRouter(db, cfg)

	staticDir := cfg.Server.StaticDir
	if info, err := os.Stat(staticDir); err == nil && info.IsDir() {
		fs := http.FileServer(http.Dir(staticDir))
		r.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"code": 4004, "msg": "not found"})
				return
			}
			fs.ServeHTTP(c.Writer, c.Request)
		})
		log.Printf("Serving frontend from %s", staticDir)
	}

	log.Printf("HU MQ started on http://localhost:%s", cfg.Server.Port)
	log.Printf("Default login: admin / admin")
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}
