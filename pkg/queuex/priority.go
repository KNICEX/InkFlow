package queuex

import (
	"cmp"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/queue"
)

type Comparable[T any] = ekit.Comparator[T]

type PriorityQueue[T any] struct {
	q *queue.PriorityQueue[T]
}

func NewPriorityQueue[T any](cap int, cmp Comparable[T]) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		q: queue.NewPriorityQueue[T](cap, ekit.Comparator[T](cmp)),
	}
}

func (q *PriorityQueue[T]) Enqueue(v T) {
	if err := q.q.Enqueue(v); err != nil {
		_, _ = q.q.Dequeue()
		_ = q.q.Enqueue(v)
	}
}

// All after call this method, the queue will be empty
func (q *PriorityQueue[T]) All() []T {
	res := make([]T, q.q.Len())
	for i := q.q.Len() - 1; i >= 0; i-- {
		t, _ := q.q.Dequeue()
		res[i] = t
	}
	return res
}

type ZQueueItem[S cmp.Ordered, V any] struct {
	Score S
	Value V
}

type ZQueue[S cmp.Ordered, V any] struct {
	q *queue.PriorityQueue[ZQueueItem[S, V]]
}

func NewZQueue[S cmp.Ordered, V any](cap int) *ZQueue[S, V] {
	return &ZQueue[S, V]{
		q: queue.NewPriorityQueue[ZQueueItem[S, V]](cap, func(a, b ZQueueItem[S, V]) int {
			if a.Score == b.Score {
				return 0
			}
			if a.Score < b.Score {
				return -1
			}
			return 1
		}),
	}
}

func (q *ZQueue[S, V]) Enqueue(score S, value V) {
	if err := q.q.Enqueue(ZQueueItem[S, V]{Score: score, Value: value}); err != nil {
		_, _ = q.q.Dequeue()
		_ = q.q.Enqueue(ZQueueItem[S, V]{Score: score, Value: value})
	}
}

func (q *ZQueue[S, V]) All() []ZQueueItem[S, V] {
	res := make([]ZQueueItem[S, V], q.q.Len())
	for i := q.q.Len() - 1; i >= 0; i-- {
		t, _ := q.q.Dequeue()
		res[i] = t
	}
	return res
}

func (q *ZQueue[S, V]) AllValues() []V {
	res := make([]V, q.q.Len())
	for i := q.q.Len() - 1; i >= 0; i-- {
		t, _ := q.q.Dequeue()
		res[i] = t.Value
	}
	return res
}
