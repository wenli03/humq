package kafka

import (
	"fmt"
	"sync"

	"github.com/IBM/sarama"
)

type AdminClient struct {
	client  sarama.Client
	admin   sarama.ClusterAdmin
	brokers string
	mu      sync.RWMutex
}

type TopicDetail struct {
	Name              string            `json:"name"`
	Partitions        int               `json:"partitions"`
	ReplicationFactor int               `json:"replication_factor"`
	PartitionDetails  []PartitionDetail `json:"partition_details"`
	Configs           map[string]string `json:"configs"`
}

type PartitionDetail struct {
	Partition int32   `json:"partition"`
	Leader    int32   `json:"leader"`
	Replicas  []int32 `json:"replicas"`
	Isr       []int32 `json:"isr"`
}

type ConsumerGroupDetail struct {
	GroupID         string           `json:"group_id"`
	Topics          []string         `json:"topics"`
	Members         int              `json:"members"`
	State           string           `json:"state"`
	Lag             int64            `json:"lag"`
	LagPerPartition map[string]int64 `json:"lag_per_partition"`
}

type ClusterInfo struct {
	Brokers        []int32 `json:"brokers"`
	TopicCount     int     `json:"topic_count"`
	PartitionCount int     `json:"partition_count"`
}

func NewAdminClient(brokers string) (*AdminClient, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0

	client, err := sarama.NewClient([]string{brokers}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}

	return &AdminClient{
		client:  client,
		admin:   admin,
		brokers: brokers,
	}, nil
}

func (a *AdminClient) Close() error {
	if a.admin != nil {
		a.admin.Close()
	}
	if a.client != nil {
		return a.client.Close()
	}
	return nil
}

func (a *AdminClient) ListTopics() ([]string, error) {
	topics, err := a.client.Topics()
	if err != nil {
		return nil, err
	}
	filtered := make([]string, 0)
	for _, t := range topics {
		if t != "__consumer_offsets" {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}

func (a *AdminClient) DescribeTopic(name string) (*TopicDetail, error) {
	metadata, err := a.admin.DescribeTopics([]string{name})
	if err != nil {
		return nil, err
	}
	if len(metadata) == 0 {
		return nil, fmt.Errorf("topic %s not found", name)
	}

	meta := metadata[0]
	detail := &TopicDetail{
		Name:             meta.Name,
		Partitions:       len(meta.Partitions),
		PartitionDetails: make([]PartitionDetail, 0),
	}

	if len(meta.Partitions) > 0 {
		detail.ReplicationFactor = len(meta.Partitions[0].Replicas)
		for _, p := range meta.Partitions {
			detail.PartitionDetails = append(detail.PartitionDetails, PartitionDetail{
				Partition: p.ID,
				Leader:    p.Leader,
				Replicas:  p.Replicas,
				Isr:       p.Isr,
			})
		}
	}

	configs, err := a.admin.DescribeConfig(sarama.ConfigResource{
		Type: sarama.TopicResource,
		Name: name,
	})
	if err == nil {
		detail.Configs = make(map[string]string)
		for _, entry := range configs {
			detail.Configs[entry.Name] = entry.Value
		}
	}

	return detail, nil
}

func (a *AdminClient) CreateTopic(name string, partitions int32, replicationFactor int16, configs map[string]*string) error {
	detail := &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}
	if configs != nil {
		detail.ConfigEntries = configs
	}
	return a.admin.CreateTopic(name, detail, false)
}

func (a *AdminClient) DeleteTopic(name string) error {
	return a.admin.DeleteTopic(name)
}

func (a *AdminClient) ListConsumerGroups() (map[string]string, error) {
	return a.admin.ListConsumerGroups()
}

func (a *AdminClient) DescribeConsumerGroup(groupID string) (*ConsumerGroupDetail, error) {
	groups, err := a.admin.DescribeConsumerGroups([]string{groupID})
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("consumer group %s not found", groupID)
	}

	g := groups[0]
	detail := &ConsumerGroupDetail{
		GroupID:         g.GroupId,
		State:           g.State,
		Members:         len(g.Members),
		Topics:          make([]string, 0),
		LagPerPartition: make(map[string]int64),
	}

	topicSet := make(map[string]bool)
	for _, m := range g.Members {
		meta, _ := m.GetMemberMetadata()
		for _, t := range meta.Topics {
			if !topicSet[t] {
				topicSet[t] = true
				detail.Topics = append(detail.Topics, t)
			}
		}
	}

	for _, topic := range detail.Topics {
		offsets, err := a.admin.ListConsumerGroupOffsets(groupID, map[string][]int32{
			topic: nil,
		})
		if err != nil {
			continue
		}
		for partition, offsetMeta := range offsets.Blocks[topic] {
			newestOffset, err := a.client.GetOffset(topic, partition, sarama.OffsetNewest)
			if err != nil {
				continue
			}
			lag := newestOffset - offsetMeta.Offset
			detail.Lag += lag
			detail.LagPerPartition[fmt.Sprintf("%s-%d", topic, partition)] = lag
		}
	}

	return detail, nil
}

func (a *AdminClient) GetClusterInfo() (*ClusterInfo, error) {
	brokers := a.client.Brokers()
	topics, err := a.client.Topics()
	if err != nil {
		return nil, err
	}

	info := &ClusterInfo{
		Brokers: make([]int32, 0),
	}
	for _, b := range brokers {
		info.Brokers = append(info.Brokers, b.ID())
	}
	for _, t := range topics {
		if t == "__consumer_offsets" {
			continue
		}
		info.TopicCount++
		partitions, _ := a.client.Partitions(t)
		info.PartitionCount += len(partitions)
	}

	return info, nil
}

func (a *AdminClient) FetchMessages(topic string, partition int32, offset int64) ([]*sarama.ConsumerMessage, error) {
	consumer, err := sarama.NewConsumerFromClient(a.client)
	if err != nil {
		return nil, err
	}
	defer consumer.Close()

	pc, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return nil, err
	}
	defer pc.Close()

	messages := make([]*sarama.ConsumerMessage, 0)
	for i := 0; i < 10; i++ {
		select {
		case msg := <-pc.Messages():
			if msg != nil {
				messages = append(messages, msg)
			}
		default:
			return messages, nil
		}
	}

	return messages, nil
}

func (a *AdminClient) GetTopicOffsets(topic string) (map[int32]int64, error) {
	partitions, err := a.client.Partitions(topic)
	if err != nil {
		return nil, err
	}

	offsets := make(map[int32]int64)
	for _, p := range partitions {
		offset, err := a.client.GetOffset(topic, p, sarama.OffsetNewest)
		if err != nil {
			continue
		}
		offsets[p] = offset
	}
	return offsets, nil
}

func (a *AdminClient) AlterTopicConfig(name string, configs map[string]*string) error {
	return a.admin.AlterConfig(sarama.TopicResource, name, configs, false)
}

func (a *AdminClient) ReassignPartitions(topics string) error {
	return nil
}

func (a *AdminClient) HealthCheck() error {
	_, err := a.client.Controller()
	return err
}

func (a *AdminClient) ListBrokers() []*sarama.Broker {
	return a.client.Brokers()
}
