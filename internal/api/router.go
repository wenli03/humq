package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wenli03/humq/internal/config"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "msg": "HU MQ running"})
	})

	return r
}
