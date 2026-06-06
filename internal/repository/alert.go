package repository

import "github.com/wenli03/humq/internal/database"

type AlertRepo struct{}

func NewAlertRepo() *AlertRepo { return &AlertRepo{} }

func (r *AlertRepo) ListRules(clusterID uint) ([]database.AlertRule, error) {
	var rules []database.AlertRule
	err := database.DB.Where("cluster_id = ?", clusterID).Order("id desc").Find(&rules).Error
	return rules, err
}

func (r *AlertRepo) CreateRule(rule *database.AlertRule) error {
	return database.DB.Create(rule).Error
}

func (r *AlertRepo) UpdateRule(rule *database.AlertRule) error {
	return database.DB.Save(rule).Error
}

func (r *AlertRepo) DeleteRule(id uint) error {
	return database.DB.Delete(&database.AlertRule{}, id).Error
}

func (r *AlertRepo) FindRuleByID(id uint) (*database.AlertRule, error) {
	var rule database.AlertRule
	err := database.DB.First(&rule, id).Error
	return &rule, err
}

func (r *AlertRepo) ListEnabledRules(clusterID uint) ([]database.AlertRule, error) {
	var rules []database.AlertRule
	err := database.DB.Where("cluster_id = ? AND enabled = true", clusterID).Find(&rules).Error
	return rules, err
}

func (r *AlertRepo) CreateEvent(event *database.AlertEvent) error {
	return database.DB.Create(event).Error
}

func (r *AlertRepo) ListEvents(clusterID uint, page, pageSize int) ([]database.AlertEvent, int64, error) {
	var events []database.AlertEvent
	var total int64
	database.DB.Model(&database.AlertEvent{}).Where("cluster_id = ?", clusterID).Count(&total)
	err := database.DB.Where("cluster_id = ?", clusterID).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("triggered_at desc").Find(&events).Error
	return events, total, err
}
