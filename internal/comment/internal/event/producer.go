package event

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const (
	topicReply      = "comment-reply"
	topicLike       = "comment-like"
	topicCancelLike = "comment-cancel-like"
	topicDelete     = "comment-delete"
)

type CommentEvtProducer interface {
	ProduceReply(ctx context.Context, event ReplyEvent) error
	ProduceLike(ctx context.Context, event LikeEvent) error
	ProduceCancelLike(ctx context.Context, event LikeEvent) error
	ProduceDelete(ctx context.Context, event DeleteEvent) error
}

type KafkaCommentEvtProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaCommentEvtProducer(producer sarama.SyncProducer) CommentEvtProducer {
	return &KafkaCommentEvtProducer{
		producer: producer,
	}
}

func (p *KafkaCommentEvtProducer) ProduceReply(ctx context.Context, event ReplyEvent) error {
	return p.produce(ctx, topicReply, event)
}

func (p *KafkaCommentEvtProducer) ProduceLike(ctx context.Context, event LikeEvent) error {
	return p.produce(ctx, topicLike, event)
}

func (p *KafkaCommentEvtProducer) ProduceDelete(ctx context.Context, event DeleteEvent) error {
	return p.produce(ctx, topicDelete, event)
}

func (p *KafkaCommentEvtProducer) ProduceCancelLike(ctx context.Context, event LikeEvent) error {
	return p.produce(ctx, topicCancelLike, event)
}

func (p *KafkaCommentEvtProducer) produce(ctx context.Context, topic string, data any) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(bs),
	})
	return err
}
