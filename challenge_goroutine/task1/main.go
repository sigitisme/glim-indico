package main

import (
	"fmt"
	"sync"
)

func main() {
	data1 := []interface{}{"bisa1", "bisa2", "bisa3"}
	data2 := []interface{}{"coba1", "coba2", "coba3"}

	count := 4

	wg := &sync.WaitGroup{}
	wg.Add(count * 2)

	printFunc := func(index int, input interface{}, wg *sync.WaitGroup) {
		fmt.Println(input, index)
		wg.Done()
	}

	for i := 1; i <= count; i++ {
		go printFunc(i, data1, wg)
		go printFunc(i, data2, wg)
	}

	wg.Wait()
}
