package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(driver, dsn string) (*gorm.DB, error) {
	var d gorm.Dialector
	switch driver {
	case "postgres":
		d = postgres.Open(dsn)
	default:
		d = sqlite.Open(dsn)
	}
	db, err := gorm.Open(d, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = db
	return db, nil
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Cluster{}, &Topic{}, &ConsumerGroup{}, &AlertRule{}, &AlertEvent{})
}

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"size:64;uniqueIndex" json:"username"`
	Password string `gorm:"size:256" json:"-"`
	Role     string `gorm:"size:16;default:user" json:"role"`
}

type Cluster struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:128" json:"name"`
	Host string `gorm:"size:256" json:"host"`
}

type Topic struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"size:256" json:"name"`
	Partitions int    `json:"partitions"`
	Replicas   int    `json:"replicas"`
}

type ConsumerGroup struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Group string `gorm:"size:256" json:"group"`
	Lag   int64  `json:"lag"`
}

type AlertRule struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `json:"name"`
	Metric    string  `json:"metric"`
	Operator  string  `json:"operator"`
	Threshold float64 `json:"threshold"`
	Enabled   bool    `json:"enabled"`
}

type AlertEvent struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

func SeedDefaultAdmin(db *gorm.DB) {
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		db.Create(&User{Username: "admin", Password: "admin", Role: "admin"})
	}
}
