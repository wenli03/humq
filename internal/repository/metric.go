package repository

import (
	"time"

	"github.com/wenli03/humq/internal/database"
)

type MetricRepo struct{}

func NewMetricRepo() *MetricRepo { return &MetricRepo{} }

func (r *MetricRepo) Insert(snapshot *database.MetricSnapshot) error {
	return database.DB.Create(snapshot).Error
}

func (r *MetricRepo) Query(clusterID uint, metricType string, start, end time.Time) ([]database.MetricSnapshot, error) {
	var snapshots []database.MetricSnapshot
	q := database.DB.Where("cluster_id = ?", clusterID)
	if metricType != "" {
		q = q.Where("metric_type = ?", metricType)
	}
	if !start.IsZero() {
		q = q.Where("collected_at >= ?", start)
	}
	if !end.IsZero() {
		q = q.Where("collected_at <= ?", end)
	}
	err := q.Order("collected_at asc").Find(&snapshots).Error
	return snapshots, err
}

func (r *MetricRepo) CleanupOlderThan(before time.Time) error {
	return database.DB.Where("collected_at < ?", before).Delete(&database.MetricSnapshot{}).Error
}
