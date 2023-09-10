package main

import (
	"fmt"
	"sync"
)

func main() {
	data1 := []interface{}{"bisa1", "bisa2", "bisa3"}
	data2 := []interface{}{"coba1", "coba2", "coba3"}

	count := 4

	mu := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// printFunc := func(index int, input interface{}, wg *sync.WaitGroup, mu *sync.Mutex) {
	// 	mu.Lock()
	// 	fmt.Println(input, index)
	// 	mu.Unlock()
	// 	wg.Done()
	// }

	go func() {
		for i := 1; i <= count; i++ {
			mu.Lock()
			fmt.Println(data1, i)
			mu.Unlock()
		}

		wg.Done()
	}()

	go func() {
		for i := 1; i <= count; i++ {
			mu.Lock()
			fmt.Println(data2, i)
			mu.Unlock()
		}

		wg.Done()
	}()

	// for i := 1; i <= count; i++ {
	// 	go printFunc(i, data1, wg, mu)
	// 	go printFunc(i, data2, wg, mu)
	// }

	wg.Wait()
}
