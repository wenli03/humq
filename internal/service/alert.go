package service

import (
	"github.com/wenli03/humq/internal/database"
	"github.com/wenli03/humq/internal/repository"
)

type AlertService struct {
	repo *repository.AlertRepo
}

func NewAlertService() *AlertService {
	return &AlertService{repo: repository.NewAlertRepo()}
}

func (s *AlertService) ListRules(clusterID uint) ([]database.AlertRule, error) {
	return s.repo.ListRules(clusterID)
}

func (s *AlertService) CreateRule(clusterID uint, name, metric, operator string, threshold float64, channels string) (*database.AlertRule, error) {
	rule := &database.AlertRule{
		ClusterID: clusterID,
		Name:      name,
		Metric:    metric,
		Operator:  operator,
		Threshold: threshold,
		Channels:  channels,
		Enabled:   true,
	}
	if err := s.repo.CreateRule(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *AlertService) UpdateRule(id uint, enabled *bool, threshold *float64) error {
	rule, err := s.repo.FindRuleByID(id)
	if err != nil {
		return err
	}
	if enabled != nil {
		rule.Enabled = *enabled
	}
	if threshold != nil {
		rule.Threshold = *threshold
	}
	return s.repo.UpdateRule(rule)
}

func (s *AlertService) DeleteRule(id uint) error {
	return s.repo.DeleteRule(id)
}

func (s *AlertService) ListEvents(clusterID uint, page, pageSize int) ([]database.AlertEvent, int64, error) {
	return s.repo.ListEvents(clusterID, page, pageSize)
}
