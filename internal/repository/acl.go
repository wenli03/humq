package repository

import "github.com/wenli03/humq/internal/database"

type ACLRepo struct{}

func NewACLRepo() *ACLRepo { return &ACLRepo{} }

func (r *ACLRepo) List() ([]database.ACLRule, error) {
	var rules []database.ACLRule
	err := database.DB.Order("id desc").Find(&rules).Error
	return rules, err
}

func (r *ACLRepo) Create(rule *database.ACLRule) error {
	return database.DB.Create(rule).Error
}

func (r *ACLRepo) Delete(id uint) error {
	return database.DB.Delete(&database.ACLRule{}, id).Error
}

func (r *ACLRepo) FindByUserAndResource(userID uint, resourceType, resourceName, operation string) (*database.ACLRule, error) {
	var rule database.ACLRule
	err := database.DB.Where("user_id = ? AND resource_type = ? AND resource_name = ? AND operation = ?",
		userID, resourceType, resourceName, operation).First(&rule).Error
	return &rule, err
}
