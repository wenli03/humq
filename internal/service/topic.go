package service

import (
	"github.com/wenli03/humq/internal/database"
	"github.com/wenli03/humq/internal/kafka"
	"github.com/wenli03/humq/internal/repository"
)

type TopicService struct {
	repo       *repository.TopicRepo
	clusterSvc *ClusterService
}

func NewTopicService(clusterSvc *ClusterService) *TopicService {
	return &TopicService{repo: repository.NewTopicRepo(), clusterSvc: clusterSvc}
}

func (s *TopicService) List(clusterID uint, keyword string) ([]database.TopicMeta, error) {
	return s.repo.ListByCluster(clusterID, keyword)
}

func (s *TopicService) Create(clusterID uint, name string, partitions int32, replicationFactor int16, configs map[string]*string) error {
	client, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}

	if err := client.CreateTopic(name, partitions, replicationFactor, configs); err != nil {
		return err
	}

	topic := &database.TopicMeta{
		ClusterID:         clusterID,
		Name:              name,
		Partitions:        int(partitions),
		ReplicationFactor: int(replicationFactor),
	}
	return s.repo.Create(topic)
}

func (s *TopicService) Delete(clusterID uint, name string) error {
	client, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	if err := client.DeleteTopic(name); err != nil {
		return err
	}
	return s.repo.Delete(clusterID, name)
}

func (s *TopicService) Describe(clusterID uint, name string) (*kafka.TopicDetail, error) {
	client, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	return client.DescribeTopic(name)
}

func (s *TopicService) AlterConfig(clusterID uint, name string, configs map[string]*string) error {
	client, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	return client.AlterTopicConfig(name, configs)
}
