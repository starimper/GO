package main

import (
	"fmt"
	"sync"
)

func main() {
	var sm sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			sm.Store("key", key)
		}(i)
	}

	wg.Wait()

	value, ok := sm.Load("key")
	if ok {
		fmt.Printf("Value: %v\n", value)
	}
}
