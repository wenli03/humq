package repository

import "github.com/wenli03/humq/internal/database"

type ConsumerRepo struct{}

func NewConsumerRepo() *ConsumerRepo { return &ConsumerRepo{} }

func (r *ConsumerRepo) ListByCluster(clusterID uint) ([]database.ConsumerGroup, error) {
	var groups []database.ConsumerGroup
	err := database.DB.Where("cluster_id = ?", clusterID).Order("id desc").Find(&groups).Error
	return groups, err
}

func (r *ConsumerRepo) FindByGroupID(clusterID uint, groupID string) (*database.ConsumerGroup, error) {
	var group database.ConsumerGroup
	err := database.DB.Where("cluster_id = ? AND group_id = ?", clusterID, groupID).First(&group).Error
	return &group, err
}

func (r *ConsumerRepo) Upsert(group *database.ConsumerGroup) error {
	var existing database.ConsumerGroup
	err := database.DB.Where("cluster_id = ? AND group_id = ?", group.ClusterID, group.GroupID).First(&existing).Error
	if err == nil {
		group.ID = existing.ID
		return database.DB.Save(group).Error
	}
	return database.DB.Create(group).Error
}
