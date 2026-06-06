package repository

import "github.com/wenli03/humq/internal/database"

type TopicRepo struct{}

func NewTopicRepo() *TopicRepo { return &TopicRepo{} }

func (r *TopicRepo) ListByCluster(clusterID uint, keyword string) ([]database.TopicMeta, error) {
	var topics []database.TopicMeta
	q := database.DB.Where("cluster_id = ?", clusterID)
	if keyword != "" {
		q = q.Where("name ILIKE ?", "%"+keyword+"%")
	}
	err := q.Order("id desc").Find(&topics).Error
	return topics, err
}

func (r *TopicRepo) FindByName(clusterID uint, name string) (*database.TopicMeta, error) {
	var topic database.TopicMeta
	err := database.DB.Where("cluster_id = ? AND name = ?", clusterID, name).First(&topic).Error
	return &topic, err
}

func (r *TopicRepo) Create(topic *database.TopicMeta) error {
	return database.DB.Create(topic).Error
}

func (r *TopicRepo) Delete(clusterID uint, name string) error {
	return database.DB.Where("cluster_id = ? AND name = ?", clusterID, name).Delete(&database.TopicMeta{}).Error
}

func (r *TopicRepo) Update(topic *database.TopicMeta) error {
	return database.DB.Save(topic).Error
}
