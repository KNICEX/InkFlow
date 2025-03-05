package event

import "context"

type InkPublishedProducer interface {
	Produce(ctx context.Context, evt InkPublishedEvt) error
}
