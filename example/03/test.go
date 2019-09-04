package main

import (
	"fmt"
	"sync"
)

var s []int32

var lock sync.Mutex

// var lock = &sync.Mutex{}

func addValue(i int32) {
	// lock.Lock()
	// defer lock.Unlock()
	s = append(s, i)
}

func main() {
	var i int32
	wg := sync.WaitGroup{}
	// defer wg.Wait()
	for i = 0; i <= 1000000; i++ {
		wg.Add(1)
		go func(i int32) {
			lock.Lock()
			addValue(i)
			defer wg.Done()
			lock.Unlock()
		}(i)
	}
	// time.Sleep(time.Second)
	wg.Wait()
	defer fmt.Println(len(s))
}
