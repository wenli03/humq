package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wenli03/humq/internal/auth"
	"github.com/wenli03/humq/internal/config"
	"github.com/wenli03/humq/internal/service"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	auth.Init(cfg.JWT.Secret, cfg.JWT.ExpireHours, cfg.JWT.RefreshExpireHours)

	r := gin.Default()

	clusterSvc := service.NewClusterService()
	topicSvc := service.NewTopicService(clusterSvc)
	consumerSvc := service.NewConsumerService(clusterSvc)
	messageSvc := service.NewMessageService(clusterSvc)
	alertSvc := service.NewAlertService()
	aclSvc := service.NewACLService()
	userSvc := service.NewUserService()

	api := r.Group("/api/v1")

	api.GET("/health", func(c *gin.Context) {
		OK(c, gin.H{"status": "running", "name": "HU MQ"})
	})

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/refresh", auth.RefreshToken)
	}

	protected := api.Group("")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/me", GetCurrentUser())

		protected.GET("/clusters", ListClusters(clusterSvc))
		protected.POST("/clusters", CreateCluster(clusterSvc))
		protected.GET("/clusters/:id/metrics", GetClusterInfo(clusterSvc))
		protected.DELETE("/clusters/:id", auth.AdminOnly(), DeleteCluster(clusterSvc))

		protected.GET("/topics", ListTopics(topicSvc))
		protected.POST("/topics", CreateTopic(topicSvc))
		protected.GET("/topics/:name/partitions", DescribeTopic(topicSvc))
		protected.PUT("/topics/:name/config", AlterTopicConfig(topicSvc))
		protected.DELETE("/topics/:name", DeleteTopic(topicSvc))

		protected.GET("/consumers", ListConsumers(consumerSvc))
		protected.GET("/consumers/:name/lag", DescribeConsumer(consumerSvc))

		protected.POST("/messages/trace", TraceMessages(messageSvc))
		protected.GET("/messages/dead", ListDeadMessages(messageSvc))
		protected.POST("/messages/replay", ReplayMessages(messageSvc))

		protected.GET("/alerts/rules", ListAlertRules(alertSvc))
		protected.POST("/alerts/rules", CreateAlertRule(alertSvc))
		protected.PUT("/alerts/rules/:id", UpdateAlertRule(alertSvc))
		protected.DELETE("/alerts/rules/:id", auth.AdminOnly(), DeleteAlertRule(alertSvc))
		protected.GET("/alerts/events", ListAlertEvents(alertSvc))

		protected.GET("/users", auth.AdminOnly(), ListUsers(userSvc))
		protected.POST("/users", auth.AdminOnly(), CreateUser(userSvc))
		protected.PUT("/users/:id", auth.AdminOnly(), UpdateUser(userSvc))
		protected.DELETE("/users/:id", auth.AdminOnly(), DeleteUser(userSvc))
		protected.GET("/acls", ListACLs(aclSvc))
		protected.POST("/acls", auth.AdminOnly(), CreateACL(aclSvc))
		protected.DELETE("/acls/:id", auth.AdminOnly(), DeleteACL(aclSvc))
	}

	return r
}
