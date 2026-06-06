package database

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedDefaultAdmin(db *gorm.DB) {
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	admin := &User{
		Username: "admin",
		Password: string(hash),
		Role:     "admin",
	}
	db.Create(admin)
}
