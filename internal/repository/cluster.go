package repository

import "github.com/wenli03/humq/internal/database"

type ClusterRepo struct{}

func NewClusterRepo() *ClusterRepo { return &ClusterRepo{} }

func (r *ClusterRepo) List() ([]database.Cluster, error) {
	var clusters []database.Cluster
	err := database.DB.Order("id desc").Find(&clusters).Error
	return clusters, err
}

func (r *ClusterRepo) Create(cluster *database.Cluster) error {
	return database.DB.Create(cluster).Error
}

func (r *ClusterRepo) FindByID(id uint) (*database.Cluster, error) {
	var cluster database.Cluster
	err := database.DB.First(&cluster, id).Error
	return &cluster, err
}

func (r *ClusterRepo) Update(cluster *database.Cluster) error {
	return database.DB.Save(cluster).Error
}

func (r *ClusterRepo) Delete(id uint) error {
	return database.DB.Delete(&database.Cluster{}, id).Error
}
