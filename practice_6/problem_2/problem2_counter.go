package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func withMutex() {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("[Mutex] Counter:", counter)
}

func withAtomic() {
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()
	fmt.Println("[Atomic] Counter:", atomic.LoadInt64(&counter))
}

func main() {
	withMutex()
	withAtomic()
}
