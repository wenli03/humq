package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wenli03/humq/internal/config"
	"gorm.io/gorm"
)

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

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": data})
}
func OKPage(c *gin.Context, data interface{}, total int64) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": data, "total": total})
}

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	api := r.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) {
		OK(c, gin.H{"status": "running", "name": "HU MQ"})
	})
	api.POST("/auth/login", handleLogin(db))

	clusters := &clusterStore{mu: sync.Mutex{}, items: []gin.H{
		{"id": 1, "name": "默认集群", "host": "localhost:9092", "status": "online", "created_at": time.Now().Format(time.RFC3339)},
	}}
	topics := &topicStore{mu: sync.Mutex{}, items: []gin.H{
		{"id": 1, "name": "order-events", "partitions": 3, "replicas": 2, "created_at": time.Now().Format(time.RFC3339)},
		{"id": 2, "name": "user-events", "partitions": 2, "replicas": 2, "created_at": time.Now().Format(time.RFC3339)},
	}}
	consumers := []gin.H{
		{"group": "order-group", "topics": []string{"order-events"}, "members": 2, "state": "Stable", "lag": 15200,
			"per_partition": []gin.H{{"partition": 0, "lag": 5200}, {"partition": 1, "lag": 4800}, {"partition": 2, "lag": 5200}}},
		{"group": "user-sync", "topics": []string{"user-events"}, "members": 1, "state": "Stable", "lag": 850,
			"per_partition": []gin.H{{"partition": 0, "lag": 450}, {"partition": 1, "lag": 400}}},
		{"group": "audit-logger", "topics": []string{"order-events", "user-events"}, "members": 2, "state": "Stable", "lag": 3200,
			"per_partition": []gin.H{{"partition": 0, "lag": 1200}, {"partition": 1, "lag": 2000}}},
	}
	alerts := []gin.H{
		{"id": 1, "name": "消费积压告警", "metric": "lag", "op": ">", "value": 10000, "enabled": true},
		{"id": 2, "name": "磁盘告警", "metric": "disk", "op": ">", "value": 80, "enabled": true},
	}
	alertEvents := []gin.H{
		{"id": 1, "level": "warning", "msg": "order-group 积压 15200 条", "status": "triggered", "time": time.Now().Format(time.RFC3339)},
	}
	users := []gin.H{
		{"id": 1, "username": "admin", "role": "admin", "created_at": "2026-06-01T00:00:00Z"},
		{"id": 2, "username": "dev", "role": "user", "created_at": "2026-06-02T00:00:00Z"},
	}
	acls := []gin.H{
		{"id": 1, "user_id": 2, "resource": "topic", "name": "order-events", "op": "read"},
	}
	backlogItems := []gin.H{
		{"consumer_group": "order-group", "topic": "order-events", "lag": 15200, "severity": "critical", "suggestion": "增加消费者至 6 个"},
		{"consumer_group": "audit-logger", "topic": "order-events", "lag": 3200, "severity": "warning", "suggestion": "检查消费者处理性能"},
	}
	opsLog := &logStore{mu: sync.Mutex{}, items: []gin.H{}}

	api.GET("/clusters", func(c *gin.Context) { OK(c, clusters.items) })
	api.POST("/clusters", func(c *gin.Context) { var req struct{ Name, Host string }; c.ShouldBindJSON(&req); clusters.add(req.Name, req.Host); OK(c, gin.H{"status": "ok"}) })
	api.GET("/clusters/:id", func(c *gin.Context) { OK(c, gin.H{"brokers": []int32{1, 2, 3}, "topics": 4, "partitions": 12}) })
	api.DELETE("/clusters/:id", func(c *gin.Context) { OK(c, nil) })

	api.GET("/topics", func(c *gin.Context) { OK(c, topics.items) })
	api.POST("/topics", func(c *gin.Context) { var req struct{ Name string; Partitions int }; c.ShouldBindJSON(&req); topics.add(req.Name, req.Partitions); OK(c, gin.H{"status": "ok"}) })
	api.GET("/topics/:name", func(c *gin.Context) {
		OK(c, gin.H{"name": c.Param("name"), "partitions": 3, "replicas": 2,
			"details": []gin.H{
				{"id": 0, "leader": 1, "replicas": []int32{1, 2}, "isr": []int32{1, 2}},
				{"id": 1, "leader": 2, "replicas": []int32{2, 3}, "isr": []int32{2, 3}},
				{"id": 2, "leader": 3, "replicas": []int32{3, 1}, "isr": []int32{3, 1}},
			},
		})
	})
	api.PUT("/topics/:name/config", func(c *gin.Context) { OK(c, gin.H{"status": "ok"}) })
	api.DELETE("/topics/:name", func(c *gin.Context) { OK(c, nil) })

	api.GET("/consumers", func(c *gin.Context) { OK(c, consumers) })
	api.GET("/consumers/:name", func(c *gin.Context) {
		for _, cg := range consumers {
			if cg["group"] == c.Param("name") {
				OK(c, cg); return
			}
		}
		OK(c, gin.H{"group": c.Param("name"), "topics": []string{}, "members": 0, "state": "Unknown", "lag": 0})
	})

	api.POST("/messages/trace", func(c *gin.Context) {
		OK(c, []gin.H{
			{"topic": "order-events", "partition": 0, "offset": 100, "key": "ord-001", "value": "{\"id\":1,\"amt\":99.9}", "time": time.Now().Format(time.RFC3339)},
			{"topic": "order-events", "partition": 0, "offset": 101, "key": "ord-002", "value": "{\"id\":2,\"amt\":150.0}", "time": time.Now().Format(time.RFC3339)},
		})
	})
	api.GET("/messages/dead", func(c *gin.Context) { OKPage(c, []gin.H{}, 0) })
	api.POST("/messages/replay", func(c *gin.Context) { OK(c, gin.H{"msg": "重放已提交"}) })

	api.GET("/alerts/rules", func(c *gin.Context) { OK(c, alerts) })
	api.POST("/alerts/rules", func(c *gin.Context) { OK(c, gin.H{"status": "ok"}) })
	api.PUT("/alerts/rules/:id", func(c *gin.Context) { OK(c, nil) })
	api.DELETE("/alerts/rules/:id", func(c *gin.Context) { OK(c, nil) })
	api.GET("/alerts/events", func(c *gin.Context) { OKPage(c, alertEvents, int64(len(alertEvents))) })

	api.GET("/users", func(c *gin.Context) { OKPage(c, users, int64(len(users))) })
	api.POST("/users", func(c *gin.Context) { OK(c, gin.H{"status": "ok"}) })
	api.PUT("/users/:id", func(c *gin.Context) { OK(c, nil) })
	api.DELETE("/users/:id", func(c *gin.Context) { OK(c, nil) })

	api.GET("/acls", func(c *gin.Context) { OK(c, acls) })
	api.POST("/acls", func(c *gin.Context) { OK(c, gin.H{"status": "ok"}) })
	api.DELETE("/acls/:id", func(c *gin.Context) { OK(c, nil) })

	api.GET("/ops/backlog", func(c *gin.Context) { OK(c, backlogItems) })
	api.POST("/ops/resolve-backlog", func(c *gin.Context) {
		var req struct {
			ConsumerGroup string `json:"consumer_group"`
			Action        string `json:"action"`
		}
		c.ShouldBindJSON(&req)
		actions := map[string]string{
			"scale":    "建议消费者扩展至分区数",
			"skip":     "已跳过 5000 条积压消息",
			"throttle": "已暂停非关键生产者",
		}
		msg := actions[req.Action]
		if msg == "" {
			msg = "方案：1)扩容消费者 2)增加分区 3)跳过旧消息 4)暂停生产者"
		}
		opsLog.add(req.ConsumerGroup, req.Action, msg)
		OK(c, gin.H{"msg": msg})
	})
	api.GET("/ops/logs", func(c *gin.Context) { OK(c, opsLog.items) })
	api.POST("/ops/rebalance", func(c *gin.Context) { opsLog.add("system", "rebalance", "分区重分配已提交"); OK(c, gin.H{"msg": "分区重分配已提交"}) })

	return r
}

func handleLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		c.ShouldBindJSON(&req)
		OK(c, gin.H{"token": "humq-token", "user": gin.H{"id": 1, "username": "admin", "role": "admin"}})
	}
}

type clusterStore struct {
	mu    sync.Mutex
	items []gin.H
}

func (s *clusterStore) add(name, host string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, gin.H{
		"id": len(s.items) + 1, "name": name, "host": host,
		"status": "online", "created_at": time.Now().Format(time.RFC3339),
	})
}

type topicStore struct {
	mu    sync.Mutex
	items []gin.H
}

func (s *topicStore) add(name string, partitions int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if partitions == 0 {
		partitions = 1
	}
	s.items = append(s.items, gin.H{
		"id": len(s.items) + 1, "name": name, "partitions": partitions,
		"replicas": 1, "created_at": time.Now().Format(time.RFC3339),
	})
}

type logStore struct {
	mu    sync.Mutex
	items []gin.H
}

func (s *logStore) add(target, action, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append([]gin.H{{
		"time":   time.Now().Format("15:04:05"),
		"target": target,
		"action": action,
		"msg":    msg,
	}}, s.items...)
}
