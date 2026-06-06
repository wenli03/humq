package repository

import "github.com/wenli03/humq/internal/database"

type DeadMsgRepo struct{}

func NewDeadMsgRepo() *DeadMsgRepo { return &DeadMsgRepo{} }

func (r *DeadMsgRepo) List(clusterID uint, topic string, page, pageSize int) ([]database.DeadMessage, int64, error) {
	var msgs []database.DeadMessage
	var total int64
	q := database.DB.Model(&database.DeadMessage{}).Where("cluster_id = ?", clusterID)
	if topic != "" {
		q = q.Where("topic = ?", topic)
	}
	q.Count(&total)
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&msgs).Error
	return msgs, total, err
}

func (r *DeadMsgRepo) Create(msg *database.DeadMessage) error {
	return database.DB.Create(msg).Error
}

func (r *DeadMsgRepo) Delete(id uint) error {
	return database.DB.Delete(&database.DeadMessage{}, id).Error
}
