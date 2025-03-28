package queuex

import (
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
