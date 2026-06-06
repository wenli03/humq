package api

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wenli03/humq/internal/auth"
	"github.com/wenli03/humq/internal/config"
	"github.com/wenli03/humq/internal/service"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config, demoMode bool) *gin.Engine {
	auth.Init(cfg.JWT.Secret, cfg.JWT.ExpireHours, cfg.JWT.RefreshExpireHours)

	r := gin.Default()

	var clusterSvc *service.ClusterService
	var topicSvc *service.TopicService
	var consumerSvc *service.ConsumerService
	var messageSvc *service.MessageService
	var alertSvc *service.AlertService
	var aclSvc *service.ACLService
	var userSvc *service.UserService

	if demoMode {
		clusterSvc = nil
		topicSvc = nil
		consumerSvc = nil
		messageSvc = nil
		alertSvc = service.NewAlertService()
		aclSvc = service.NewACLService()
		userSvc = service.NewUserService()

		r.Use(CORSMiddleware())
	} else {
		clusterSvc = service.NewClusterService()
		topicSvc = service.NewTopicService(clusterSvc)
		consumerSvc = service.NewConsumerService(clusterSvc)
		messageSvc = service.NewMessageService(clusterSvc)
		alertSvc = service.NewAlertService()
		aclSvc = service.NewACLService()
		userSvc = service.NewUserService()
	}

	api := r.Group("/api/v1")

	api.GET("/health", func(c *gin.Context) {
		OK(c, gin.H{"status": "running", "name": "HU MQ", "demo": demoMode})
	})

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", auth.DemoLoginHandler(demoMode))
		authGroup.POST("/refresh", auth.RefreshToken)
	}

	protected := api.Group("")
	if !demoMode {
		protected.Use(auth.AuthMiddleware())
	}
	{
		protected.GET("/me", GetCurrentUser())

		if demoMode {
			protected.GET("/clusters", DemoListClusters())
			protected.POST("/clusters", DemoCreateCluster())
			protected.GET("/clusters/:id/metrics", DemoGetClusterInfo())
			protected.DELETE("/clusters/:id", DemoDeleteCluster())

			protected.GET("/topics", DemoListTopics())
			protected.POST("/topics", DemoCreateTopic())
			protected.GET("/topics/:name/partitions", DemoDescribeTopic())
			protected.PUT("/topics/:name/config", DemoAlterTopicConfig())
			protected.DELETE("/topics/:name", DemoDeleteTopic())

			protected.GET("/consumers", DemoListConsumers())
			protected.GET("/consumers/:name/lag", DemoDescribeConsumer())

			protected.POST("/messages/trace", DemoTraceMessages())
			protected.GET("/messages/dead", DemoListDeadMessages())
			protected.POST("/messages/replay", DemoReplayMessages())

			protected.GET("/alerts/rules", DemoListAlertRules())
			protected.POST("/alerts/rules", DemoCreateAlertRule())
			protected.PUT("/alerts/rules/:id", DemoUpdateAlertRule())
			protected.DELETE("/alerts/rules/:id", DemoDeleteAlertRule())
			protected.GET("/alerts/events", DemoListAlertEvents())

			protected.GET("/users", DemoListUsers())
			protected.POST("/users", DemoCreateUser())
			protected.PUT("/users/:id", DemoUpdateUser())
			protected.DELETE("/users/:id", DemoDeleteUser())
			protected.GET("/acls", DemoListACLs())
			protected.POST("/acls", DemoCreateACL())
			protected.DELETE("/acls/:id", DemoDeleteACL())

			protected.POST("/ops/rebalance", DemoRebalance())
			protected.GET("/ops/backlog", DemoListBacklog())
			protected.POST("/ops/resolve-backlog", DemoResolveBacklog())
		} else {
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
	}

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

var demoMu sync.Mutex
var demoClusters = []gin.H{
	{"id": 1, "name": "生产集群", "bootstrap_servers": "kafka-prod:9092", "status": "online", "created_at": "2026-06-06T10:00:00Z"},
	{"id": 2, "name": "测试集群", "bootstrap_servers": "kafka-test:9092", "status": "online", "created_at": "2026-06-05T08:00:00Z"},
}
var demoTopicSeq = 3
var demoTopics = []gin.H{
	{"id": 1, "cluster_id": 1, "name": "order-events", "partitions": 3, "replication_factor": 2, "created_at": "2026-06-06T10:00:00Z"},
	{"id": 2, "cluster_id": 1, "name": "user-events", "partitions": 2, "replication_factor": 2, "created_at": "2026-06-06T09:00:00Z"},
}
var demoConsumers = []gin.H{
	{"group_id": "order-processor", "topics": []string{"order-events"}, "members": 2, "state": "Stable", "lag": 15200, "lag_per_partition": gin.H{"order-events-0": 5200, "order-events-1": 4800, "order-events-2": 5200}},
	{"group_id": "user-sync", "topics": []string{"user-events"}, "members": 1, "state": "Stable", "lag": 850, "lag_per_partition": gin.H{"user-events-0": 450, "user-events-1": 400}},
	{"group_id": "log-archiver", "topics": []string{"order-events", "user-events"}, "members": 3, "state": "Stable", "lag": 2350, "lag_per_partition": gin.H{"order-events-0": 800, "user-events-0": 1550}},
}
var demoAlertRules = []gin.H{
	{"id": 1, "cluster_id": 1, "name": "消费积压告警", "metric": "lag", "operator": ">", "threshold": 10000, "enabled": true},
	{"id": 2, "cluster_id": 1, "name": "磁盘使用率告警", "metric": "disk_usage", "operator": ">", "threshold": 80, "enabled": true},
}
var demoAlertEvents = []gin.H{
	{"id": 1, "rule_id": 1, "cluster_id": 1, "level": "warning", "message": "消费组 order-processor 积压 15200 条", "status": "triggered", "triggered_at": "2026-06-06T11:00:00Z"},
}
var demoUsers = []gin.H{
	{"id": 1, "username": "admin", "role": "admin", "created_at": "2026-06-01T00:00:00Z"},
	{"id": 2, "username": "operator", "role": "user", "created_at": "2026-06-02T00:00:00Z"},
}
var demoACLs = []gin.H{
	{"id": 1, "user_id": 2, "resource_type": "topic", "resource_name": "order-events", "operation": "read"},
}
var demoBacklog = []gin.H{
	{"consumer_group": "order-processor", "topic": "order-events", "lag": 15200, "severity": "critical", "suggestion": "增加消费者实例至 6 个"},
	{"consumer_group": "log-archiver", "topic": "order-events", "lag": 2350, "severity": "warning", "suggestion": "检查消费者处理性能"},
}

func DemoListClusters() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, demoClusters) }
}
func DemoCreateCluster() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name             string `json:"name"`
			BootstrapServers string `json:"bootstrap_servers"`
		}
		c.ShouldBindJSON(&req)
		demoMu.Lock()
		demoClusters = append(demoClusters, gin.H{"id": len(demoClusters) + 1, "name": req.Name, "bootstrap_servers": req.BootstrapServers, "status": "online", "created_at": time.Now().Format(time.RFC3339)})
		demoMu.Unlock()
		OK(c, gin.H{"id": len(demoClusters)})
	}
}
func DemoGetClusterInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		OK(c, gin.H{"brokers": []int32{1, 2, 3}, "topic_count": 4, "partition_count": 12})
	}
}
func DemoDeleteCluster() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, nil) }
}

func DemoListTopics() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, demoTopics) }
}
func DemoCreateTopic() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ClusterID         uint   `json:"cluster_id"`
			Name              string `json:"name"`
			Partitions        int    `json:"partitions"`
			ReplicationFactor int    `json:"replication_factor"`
		}
		c.ShouldBindJSON(&req)
		demoMu.Lock()
		demoTopicSeq++
		t := gin.H{"id": demoTopicSeq, "cluster_id": req.ClusterID, "name": req.Name, "partitions": req.Partitions, "replication_factor": req.ReplicationFactor, "created_at": time.Now().Format(time.RFC3339)}
		demoTopics = append(demoTopics, t)
		demoMu.Unlock()
		OK(c, t)
	}
}
func DemoDescribeTopic() gin.HandlerFunc {
	return func(c *gin.Context) {
		OK(c, gin.H{
			"name": c.Param("name"), "partitions": 3, "replication_factor": 2,
			"partition_details": []gin.H{
				{"partition": 0, "leader": 1, "replicas": []int32{1, 2}, "isr": []int32{1, 2}},
				{"partition": 1, "leader": 2, "replicas": []int32{2, 3}, "isr": []int32{2, 3}},
				{"partition": 2, "leader": 3, "replicas": []int32{3, 1}, "isr": []int32{3, 1}},
			},
			"configs": gin.H{"retention.ms": "604800000"},
		})
	}
}
func DemoAlterTopicConfig() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, nil) }
}
func DemoDeleteTopic() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, nil) }
}

func DemoListConsumers() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, demoConsumers) }
}
func DemoDescribeConsumer() gin.HandlerFunc {
	return func(c *gin.Context) {
		OK(c, gin.H{
			"group_id": c.Param("name"), "topics": []string{"order-events"}, "members": 2, "state": "Stable", "lag": 15200,
			"lag_per_partition": gin.H{"order-events-0": 5200, "order-events-1": 4800, "order-events-2": 5200},
		})
	}
}

func DemoTraceMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		OK(c, []gin.H{
			{"topic": "order-events", "partition": 0, "offset": 100, "key": "order-123", "value": "{\"orderId\":123,\"amount\":99.9}", "timestamp": "2026-06-06T11:00:00Z"},
			{"topic": "order-events", "partition": 0, "offset": 101, "key": "order-124", "value": "{\"orderId\":124,\"amount\":150.0}", "timestamp": "2026-06-06T11:00:01Z"},
		})
	}
}
func DemoListDeadMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		OK(c, gin.H{"data": []gin.H{
			{"id": 1, "cluster_id": 1, "topic": "order-events", "partition": 0, "offset": 99, "key": "bad-order", "payload": "corrupted-data", "timestamp": "2026-06-06T10:00:00Z"},
		}, "total": 1})
	}
}
func DemoReplayMessages() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, gin.H{"msg": "消息重放已提交"}) }
}

func DemoListAlertRules() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, demoAlertRules) }
}
func DemoCreateAlertRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ClusterID uint    `json:"cluster_id"`
			Name      string  `json:"name"`
			Metric    string  `json:"metric"`
			Operator  string  `json:"operator"`
			Threshold float64 `json:"threshold"`
		}
		c.ShouldBindJSON(&req)
		demoMu.Lock()
		rule := gin.H{"id": len(demoAlertRules) + 1, "name": req.Name, "metric": req.Metric, "operator": req.Operator, "threshold": req.Threshold, "enabled": true}
		demoAlertRules = append(demoAlertRules, rule)
		demoMu.Unlock()
		OK(c, rule)
	}
}
func DemoUpdateAlertRule() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, nil) }
}
func DemoDeleteAlertRule() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, nil) }
}
func DemoListAlertEvents() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, gin.H{"data": demoAlertEvents, "total": len(demoAlertEvents)}) }
}

func DemoListUsers() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, gin.H{"data": demoUsers, "total": len(demoUsers)}) }
}
func DemoCreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		c.ShouldBindJSON(&req)
		demoMu.Lock()
		demoUsers = append(demoUsers, gin.H{"id": len(demoUsers) + 1, "username": req.Username, "role": req.Role, "created_at": time.Now().Format(time.RFC3339)})
		demoMu.Unlock()
		OK(c, gin.H{"id": len(demoUsers)})
	}
}
func DemoUpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, nil) }
}
func DemoDeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, nil) }
}

func DemoListACLs() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, demoACLs) }
}
func DemoCreateACL() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID       uint   `json:"user_id"`
			ResourceType string `json:"resource_type"`
			ResourceName string `json:"resource_name"`
			Operation    string `json:"operation"`
		}
		c.ShouldBindJSON(&req)
		demoMu.Lock()
		demoACLs = append(demoACLs, gin.H{"id": len(demoACLs) + 1, "user_id": req.UserID, "resource_type": req.ResourceType, "resource_name": req.ResourceName, "operation": req.Operation})
		demoMu.Unlock()
		OK(c, nil)
	}
}
func DemoDeleteACL() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, nil) }
}

func DemoRebalance() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, gin.H{"msg": "分区重分配已提交"}) }
}
func DemoListBacklog() gin.HandlerFunc {
	return func(c *gin.Context) { OK(c, demoBacklog) }
}
func DemoResolveBacklog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ConsumerGroup string `json:"consumer_group"`
			Action        string `json:"action"`
		}
		c.ShouldBindJSON(&req)
		if req.Action == "scale" {
			OK(c, gin.H{"msg": "建议扩容消费者：增加消费者实例数 = 原数量 x 2"})
		} else if req.Action == "skip" {
			OK(c, gin.H{"msg": "已跳过最旧的 5000 条消息，积压降低"})
		} else if req.Action == "throttle" {
			OK(c, gin.H{"msg": "已暂停非关键生产者"})
		} else {
			OK(c, gin.H{"msg": "积压处理方案：1.扩容消费者 2.增加分区 3.跳过旧消息 4.暂停生产者"})
		}
	}
}
