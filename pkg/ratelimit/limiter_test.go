package ratelimit

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestKeyWindowSliceLimiter_Limited(t *testing.T) {
	lim := NewSliceWindowKeyLimiter(time.Second, 400)
	wg := &sync.WaitGroup{}
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 300; j++ {
				ok, err := lim.Limited(nil, fmt.Sprintf("test_%d", i))
				if err != nil {
					t.Error(err)
				} else if ok {
					t.Log("限流了")
				} else {
					t.Log("没有限流")
				}
			}
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second * 2)
}
