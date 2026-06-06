package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wenli03/humq/internal/service"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Msg: "success", Data: data})
}

func OKPage(c *gin.Context, data interface{}, total int64) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": data, "total": total})
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{Code: code, Msg: msg, Data: nil})
}

func ListClusters(svc *service.ClusterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusters, err := svc.List()
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, clusters)
	}
}

func CreateCluster(svc *service.ClusterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name             string `json:"name" binding:"required"`
			BootstrapServers string `json:"bootstrap_servers" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, 4000, "参数错误")
			return
		}
		cluster, err := svc.Create(req.Name, req.BootstrapServers)
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, cluster)
	}
}

func DeleteCluster(svc *service.ClusterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := svc.Delete(uint(id)); err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, nil)
	}
}

func GetClusterInfo(svc *service.ClusterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		info, err := svc.GetInfo(uint(id))
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, info)
	}
}

func ListTopics(svc *service.TopicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID, _ := strconv.Atoi(c.Query("cluster_id"))
		keyword := c.Query("keyword")
		topics, err := svc.List(uint(clusterID), keyword)
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, topics)
	}
}

func CreateTopic(svc *service.TopicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ClusterID         uint              `json:"cluster_id" binding:"required"`
			Name              string            `json:"name" binding:"required"`
			Partitions        int32             `json:"partitions"`
			ReplicationFactor int16             `json:"replication_factor"`
			Configs           map[string]string `json:"configs"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, 4000, "参数错误")
			return
		}
		if req.Partitions == 0 {
			req.Partitions = 1
		}
		if req.ReplicationFactor == 0 {
			req.ReplicationFactor = 1
		}
		configPtrs := make(map[string]*string)
		for k, v := range req.Configs {
			val := v
			configPtrs[k] = &val
		}
		if err := svc.Create(req.ClusterID, req.Name, req.Partitions, req.ReplicationFactor, configPtrs); err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, nil)
	}
}

func DeleteTopic(svc *service.TopicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID, _ := strconv.Atoi(c.Query("cluster_id"))
		name := c.Param("name")
		if err := svc.Delete(uint(clusterID), name); err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, nil)
	}
}

func DescribeTopic(svc *service.TopicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID, _ := strconv.Atoi(c.Query("cluster_id"))
		name := c.Param("name")
		detail, err := svc.Describe(uint(clusterID), name)
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, detail)
	}
}

func AlterTopicConfig(svc *service.TopicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ClusterID uint              `json:"cluster_id" binding:"required"`
			Configs   map[string]string `json:"configs" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, 4000, "参数错误")
			return
		}
		name := c.Param("name")
		configPtrs := make(map[string]*string)
		for k, v := range req.Configs {
			val := v
			configPtrs[k] = &val
		}
		if err := svc.AlterConfig(req.ClusterID, name, configPtrs); err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, nil)
	}
}

func ListConsumers(svc *service.ConsumerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID, _ := strconv.Atoi(c.Query("cluster_id"))
		groups, err := svc.List(uint(clusterID))
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, groups)
	}
}

func DescribeConsumer(svc *service.ConsumerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID, _ := strconv.Atoi(c.Query("cluster_id"))
		name := c.Param("name")
		detail, err := svc.Describe(uint(clusterID), name)
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, detail)
	}
}

func TraceMessages(svc *service.MessageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ClusterID uint   `json:"cluster_id" binding:"required"`
			Topic     string `json:"topic" binding:"required"`
			Partition int32  `json:"partition"`
			Offset    int64  `json:"offset"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, 4000, "参数错误")
			return
		}
		msgs, err := svc.Trace(req.ClusterID, req.Topic, req.Partition, req.Offset)
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, msgs)
	}
}

func ListDeadMessages(svc *service.MessageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID, _ := strconv.Atoi(c.Query("cluster_id"))
		topic := c.Query("topic")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		msgs, total, err := svc.ListDead(uint(clusterID), topic, page, pageSize)
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OKPage(c, msgs, total)
	}
}

func ReplayMessages(svc *service.MessageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ClusterID   uint   `json:"cluster_id" binding:"required"`
			DeadMsgID   uint   `json:"dead_msg_id" binding:"required"`
			TargetTopic string `json:"target_topic" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, 4000, "参数错误")
			return
		}
		if err := svc.Replay(req.ClusterID, req.DeadMsgID, req.TargetTopic); err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, nil)
	}
}

func ListAlertRules(svc *service.AlertService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID, _ := strconv.Atoi(c.Query("cluster_id"))
		rules, err := svc.ListRules(uint(clusterID))
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, rules)
	}
}

func CreateAlertRule(svc *service.AlertService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ClusterID uint    `json:"cluster_id" binding:"required"`
			Name      string  `json:"name" binding:"required"`
			Metric    string  `json:"metric" binding:"required"`
			Operator  string  `json:"operator" binding:"required"`
			Threshold float64 `json:"threshold" binding:"required"`
			Channels  string  `json:"channels"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, 4000, "参数错误")
			return
		}
		rule, err := svc.CreateRule(req.ClusterID, req.Name, req.Metric, req.Operator, req.Threshold, req.Channels)
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, rule)
	}
}

func UpdateAlertRule(svc *service.AlertService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var req struct {
			Enabled   *bool    `json:"enabled"`
			Threshold *float64 `json:"threshold"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, 4000, "参数错误")
			return
		}
		if err := svc.UpdateRule(uint(id), req.Enabled, req.Threshold); err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, nil)
	}
}

func DeleteAlertRule(svc *service.AlertService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := svc.DeleteRule(uint(id)); err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, nil)
	}
}

func ListAlertEvents(svc *service.AlertService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID, _ := strconv.Atoi(c.Query("cluster_id"))
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		events, total, err := svc.ListEvents(uint(clusterID), page, pageSize)
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OKPage(c, events, total)
	}
}

func ListUsers(svc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		users, total, err := svc.List(page, pageSize)
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OKPage(c, users, total)
	}
}

func CreateUser(svc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
			Role     string `json:"role"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, 4000, "参数错误")
			return
		}
		user, err := svc.Create(req.Username, req.Password, req.Role)
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, user)
	}
}

func UpdateUser(svc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var req struct {
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, 4000, "参数错误")
			return
		}
		if err := svc.Update(uint(id), req.Password, req.Role); err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, nil)
	}
}

func DeleteUser(svc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := svc.Delete(uint(id)); err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, nil)
	}
}

func ListACLs(svc *service.ACLService) gin.HandlerFunc {
	return func(c *gin.Context) {
		rules, err := svc.List()
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, rules)
	}
}

func CreateACL(svc *service.ACLService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID       uint   `json:"user_id" binding:"required"`
			ResourceType string `json:"resource_type" binding:"required"`
			ResourceName string `json:"resource_name" binding:"required"`
			Operation    string `json:"operation" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, 4000, "参数错误")
			return
		}
		rule, err := svc.Create(req.UserID, req.ResourceType, req.ResourceName, req.Operation)
		if err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, rule)
	}
}

func DeleteACL(svc *service.ACLService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := svc.Delete(uint(id)); err != nil {
			Fail(c, 5000, err.Error())
			return
		}
		OK(c, nil)
	}
}

func GetCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		role, _ := c.Get("role")
		OK(c, gin.H{"id": userID, "username": username, "role": role})
	}
}
