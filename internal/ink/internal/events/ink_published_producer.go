package events

import "context"

type InkPublishedProducer interface {
	Produce(ctx context.Context, evt InkPublishedEvt) error
}
