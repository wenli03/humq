package repository

import (
	"github.com/wenli03/humq/internal/database"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo { return &UserRepo{} }

func (r *UserRepo) FindByUsername(username string) (*database.User, error) {
	var user database.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByID(id uint) (*database.User, error) {
	var user database.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) List(page, pageSize int) ([]database.User, int64, error) {
	var users []database.User
	var total int64
	database.DB.Model(&database.User{}).Count(&total)
	err := database.DB.Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&users).Error
	return users, total, err
}

func (r *UserRepo) Create(user *database.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepo) Update(user *database.User) error {
	return database.DB.Save(user).Error
}

func (r *UserRepo) Delete(id uint) error {
	return database.DB.Delete(&database.User{}, id).Error
}
