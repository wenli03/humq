package service

import (
	"errors"

	"github.com/wenli03/humq/internal/database"
	"github.com/wenli03/humq/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService() *UserService {
	return &UserService{repo: repository.NewUserRepo()}
}

func (s *UserService) List(page, pageSize int) ([]database.User, int64, error) {
	return s.repo.List(page, pageSize)
}

func (s *UserService) Create(username, password, role string) (*database.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("用户名和密码不能为空")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &database.User{
		Username: username,
		Password: string(hash),
		Role:     role,
	}
	if role == "" {
		user.Role = "user"
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(id uint, password, role string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hash)
	}
	if role != "" {
		user.Role = role
	}
	return s.repo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}
