package service

import (
	"strings"

	"github.com/wenli03/humq/internal/database"
	"github.com/wenli03/humq/internal/kafka"
	"github.com/wenli03/humq/internal/repository"
)

type ConsumerService struct {
	repo       *repository.ConsumerRepo
	clusterSvc *ClusterService
}

func NewConsumerService(clusterSvc *ClusterService) *ConsumerService {
	return &ConsumerService{repo: repository.NewConsumerRepo(), clusterSvc: clusterSvc}
}

func (s *ConsumerService) List(clusterID uint) ([]kafka.ConsumerGroupDetail, error) {
	client, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}

	groups, err := client.ListConsumerGroups()
	if err != nil {
		return nil, err
	}

	details := make([]kafka.ConsumerGroupDetail, 0)
	for groupID := range groups {
		detail, err := client.DescribeConsumerGroup(groupID)
		if err != nil {
			continue
		}
		details = append(details, *detail)

		cg := &database.ConsumerGroup{
			ClusterID: clusterID,
			GroupID:   groupID,
			Topics:    strings.Join(detail.Topics, ","),
			Members:   detail.Members,
			State:     detail.State,
		}
		s.repo.Upsert(cg)
	}

	return details, nil
}

func (s *ConsumerService) Describe(clusterID uint, groupID string) (*kafka.ConsumerGroupDetail, error) {
	client, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	return client.DescribeConsumerGroup(groupID)
}
