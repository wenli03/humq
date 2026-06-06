package service

import (
	"time"

	"github.com/wenli03/humq/internal/database"
	"github.com/wenli03/humq/internal/repository"
)

type MessageService struct {
	deadMsgRepo *repository.DeadMsgRepo
	clusterSvc  *ClusterService
}

func NewMessageService(clusterSvc *ClusterService) *MessageService {
	return &MessageService{
		deadMsgRepo: repository.NewDeadMsgRepo(),
		clusterSvc:  clusterSvc,
	}
}

type TraceResult struct {
	Topic     string    `json:"topic"`
	Partition int32     `json:"partition"`
	Offset    int64     `json:"offset"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *MessageService) Trace(clusterID uint, topic string, partition int32, offset int64) ([]TraceResult, error) {
	client, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}

	msgs, err := client.FetchMessages(topic, partition, offset)
	if err != nil {
		return nil, err
	}

	results := make([]TraceResult, 0, len(msgs))
	for _, m := range msgs {
		results = append(results, TraceResult{
			Topic:     m.Topic,
			Partition: m.Partition,
			Offset:    m.Offset,
			Key:       string(m.Key),
			Value:     string(m.Value),
			Timestamp: m.Timestamp,
		})
	}

	return results, nil
}

func (s *MessageService) ListDead(clusterID uint, topic string, page, pageSize int) ([]database.DeadMessage, int64, error) {
	return s.deadMsgRepo.List(clusterID, topic, page, pageSize)
}

func (s *MessageService) Replay(clusterID uint, deadMsgID uint, targetTopic string) error {
	return nil
}
