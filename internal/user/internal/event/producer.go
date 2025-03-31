package event

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const (
	topicUserCreate = "user-create"
	topicUserUpdate = "user-update"
)

type UserProducer interface {
	ProduceCreate(ctx context.Context, event UserCreateEvent) error
	ProduceUpdate(ctx context.Context, event UserUpdateEvent) error
}

type KafkaUserProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaUserProducer(producer sarama.SyncProducer) UserProducer {
	return &KafkaUserProducer{
		producer: producer,
	}
}

func (p *KafkaUserProducer) ProduceCreate(ctx context.Context, event UserCreateEvent) error {
	bs, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicUserCreate,
		Value: sarama.ByteEncoder(bs),
	})
	return err
}

func (p *KafkaUserProducer) ProduceUpdate(ctx context.Context, event UserUpdateEvent) error {
	bs, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicUserUpdate,
		Value: sarama.ByteEncoder(bs),
	})
	return err
}
