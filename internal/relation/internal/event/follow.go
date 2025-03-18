package event

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
)

const topicFollowEvt = "user_follow_event"

type FollowEvt struct {
	FollowerId int64
	FolloweeId int64
	CreatedAt  time.Time
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
		Topic: topicFollowEvt,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
