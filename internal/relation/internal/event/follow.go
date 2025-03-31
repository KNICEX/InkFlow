package event

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
)

const topicFollow = "user-follow"

type FollowEvt struct {
	FollowerId int64     `json:"followerId"`
	FolloweeId int64     `json:"followeeId"`
	CreatedAt  time.Time `json:"createdAt"`
}

type FollowProducer interface {
	Produce(ctx context.Context, evt FollowEvt) error
}

type KafkaFollowProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaFollowProducer(producer sarama.SyncProducer) FollowProducer {
	return &KafkaFollowProducer{
		producer: producer,
	}
}

func (p *KafkaFollowProducer) Produce(ctx context.Context, evt FollowEvt) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicFollow,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
