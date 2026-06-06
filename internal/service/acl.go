package service

import (
	"github.com/wenli03/humq/internal/database"
	"github.com/wenli03/humq/internal/repository"
)

type ACLService struct {
	repo *repository.ACLRepo
}

func NewACLService() *ACLService {
	return &ACLService{repo: repository.NewACLRepo()}
}

func (s *ACLService) List() ([]database.ACLRule, error) {
	return s.repo.List()
}

func (s *ACLService) Create(userID uint, resourceType, resourceName, operation string) (*database.ACLRule, error) {
	rule := &database.ACLRule{
		UserID:       userID,
		ResourceType: resourceType,
		ResourceName: resourceName,
		Operation:    operation,
	}
	if err := s.repo.Create(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *ACLService) Delete(id uint) error {
	return s.repo.Delete(id)
}
