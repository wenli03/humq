package service

import (
	"sync"

	"github.com/wenli03/humq/internal/database"
	"github.com/wenli03/humq/internal/kafka"
	"github.com/wenli03/humq/internal/repository"
)

type ClusterService struct {
	repo    *repository.ClusterRepo
	clients map[uint]*kafka.AdminClient
	mu      sync.RWMutex
}

func NewClusterService() *ClusterService {
	return &ClusterService{
		repo:    repository.NewClusterRepo(),
		clients: make(map[uint]*kafka.AdminClient),
	}
}

func (s *ClusterService) List() ([]database.Cluster, error) {
	return s.repo.List()
}

func (s *ClusterService) Create(name, bootstrapServers string) (*database.Cluster, error) {
	cluster := &database.Cluster{
		Name:             name,
		BootstrapServers: bootstrapServers,
		Status:           "unknown",
	}
	if err := s.repo.Create(cluster); err != nil {
		return nil, err
	}
	return cluster, nil
}

func (s *ClusterService) GetClient(clusterID uint) (*kafka.AdminClient, error) {
	s.mu.RLock()
	if c, ok := s.clients[clusterID]; ok {
		s.mu.RUnlock()
		if c.HealthCheck() == nil {
			return c, nil
		}
		c.Close()
	} else {
		s.mu.RUnlock()
	}

	cluster, err := s.repo.FindByID(clusterID)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if c, ok := s.clients[clusterID]; ok {
		c.Close()
	}

	client, err := kafka.NewAdminClient(cluster.BootstrapServers)
	if err != nil {
		cluster.Status = "offline"
		s.repo.Update(cluster)
		return nil, err
	}

	cluster.Status = "online"
	s.repo.Update(cluster)
	s.clients[clusterID] = client
	return client, nil
}

func (s *ClusterService) GetInfo(clusterID uint) (*kafka.ClusterInfo, error) {
	client, err := s.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	return client.GetClusterInfo()
}

func (s *ClusterService) Delete(id uint) error {
	s.mu.Lock()
	if c, ok := s.clients[id]; ok {
		c.Close()
		delete(s.clients, id)
	}
	s.mu.Unlock()
	return s.repo.Delete(id)
}
