package queuex

import (
	"fmt"
	"testing"
)

func TestPriority(t *testing.T) {
	q := NewPriorityQueue(10, func(src int, dst int) int {
		return src - dst
	})

	for i := 0; i < 10000; i++ {
		q.Enqueue(i)
	}

	fmt.Println(q.All())
}
