package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const (
	inkViewTopic       = "ink_view"
	inkLikeTopic       = "ink_like"
	inkCancelLikeTopic = "ink_cancel_like"
)

type InteractiveProducer interface {
	ProduceInkView(ctx context.Context, evt InkViewEvent) error
	ProduceInkLike(ctx context.Context, evt InkLikeEvent) error
	ProduceInkCancelLike(ctx context.Context, evt InkCancelLikeEvent) error
}

type KafkaInteractiveProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaInteractiveProducer(producer sarama.SyncProducer) InteractiveProducer {
	return &KafkaInteractiveProducer{
		producer: producer,
	}
}

func (p *KafkaInteractiveProducer) produce(ctx context.Context, topic string, evt any) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func (p *KafkaInteractiveProducer) ProduceInkView(ctx context.Context, evt InkViewEvent) error {
	return p.produce(ctx, inkViewTopic, evt)
}

func (p *KafkaInteractiveProducer) ProduceInkLike(ctx context.Context, evt InkLikeEvent) error {
	return p.produce(ctx, inkLikeTopic, evt)
}

func (p *KafkaInteractiveProducer) ProduceInkCancelLike(ctx context.Context, evt InkCancelLikeEvent) error {
	return p.produce(ctx, inkCancelLikeTopic, evt)
}
