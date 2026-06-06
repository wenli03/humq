package database

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password  string    `gorm:"size:256;not null" json:"-"`
	Role      string    `gorm:"size:16;default:user" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Cluster struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"uniqueIndex;size:128;not null" json:"name"`
	BootstrapServers string    `gorm:"size:512;not null" json:"bootstrap_servers"`
	Status           string    `gorm:"size:16;default:unknown" json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type TopicMeta struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ClusterID         uint      `gorm:"index;not null" json:"cluster_id"`
	Name              string    `gorm:"size:256;not null" json:"name"`
	Partitions        int       `gorm:"default:1" json:"partitions"`
	ReplicationFactor int       `gorm:"default:1" json:"replication_factor"`
	RetentionMs       int64     `gorm:"default:604800000" json:"retention_ms"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ConsumerGroup struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ClusterID uint      `gorm:"index;not null" json:"cluster_id"`
	GroupID   string    `gorm:"size:256;not null" json:"group_id"`
	Topics    string    `gorm:"size:1024" json:"topics"`
	Members   int       `gorm:"default:0" json:"members"`
	State     string    `gorm:"size:32" json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MetricSnapshot struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ClusterID   uint      `gorm:"index;not null" json:"cluster_id"`
	MetricType  string    `gorm:"size:64;not null;index" json:"metric_type"`
	Value       float64   `json:"value"`
	Tags        string    `gorm:"type:jsonb;default:'{}'" json:"tags"`
	CollectedAt time.Time `gorm:"index" json:"collected_at"`
}

type AlertRule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ClusterID uint      `gorm:"index;not null" json:"cluster_id"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	Metric    string    `gorm:"size:64;not null" json:"metric"`
	Operator  string    `gorm:"size:8;not null" json:"operator"`
	Threshold float64   `json:"threshold"`
	Channels  string    `gorm:"type:jsonb;default:'[]'" json:"channels"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AlertEvent struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	RuleID      uint       `gorm:"index;not null" json:"rule_id"`
	ClusterID   uint       `gorm:"index;not null" json:"cluster_id"`
	Level       string     `gorm:"size:16;not null" json:"level"`
	Message     string     `gorm:"size:1024" json:"message"`
	Status      string     `gorm:"size:16;default:triggered" json:"status"`
	TriggeredAt time.Time  `gorm:"index" json:"triggered_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
}

type DeadMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ClusterID uint      `gorm:"index;not null" json:"cluster_id"`
	Topic     string    `gorm:"size:256;not null" json:"topic"`
	Partition int32     `json:"partition"`
	Offset    int64     `json:"offset"`
	Key       string    `gorm:"size:1024" json:"key"`
	Payload   string    `gorm:"type:text" json:"payload"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

type ACLRule struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	ResourceType string    `gorm:"size:32;not null" json:"resource_type"`
	ResourceName string    `gorm:"size:256;not null" json:"resource_name"`
	Operation    string    `gorm:"size:32;not null" json:"operation"`
	CreatedAt    time.Time `json:"created_at"`
}

func (User) TableName() string          { return "users" }
func (Cluster) TableName() string       { return "clusters" }
func (TopicMeta) TableName() string     { return "topics_meta" }
func (ConsumerGroup) TableName() string  { return "consumer_groups" }
func (MetricSnapshot) TableName() string { return "metrics_snapshots" }
func (AlertRule) TableName() string      { return "alert_rules" }
func (AlertEvent) TableName() string     { return "alert_events" }
func (DeadMessage) TableName() string    { return "dead_messages" }
func (ACLRule) TableName() string        { return "acl_rules" }
