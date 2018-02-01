package main

import (
	"fmt"
	"sync"
)

type pair struct {
	x int
	y int
}

func main() {
	pairs := []*pair{}
	pairChan := make(chan *pair)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			p := handleNumber(val)
			fmt.Printf("%+v\n", p)
			pairChan <- p
			wg.Done()
		}(i)
	}
	go func() {
		for p := range pairChan {
			pairs = append(pairs, p)
		}
	}()
	wg.Wait()
	close(pairChan)
	fmt.Println("Done")
}

func handleNumber(i int) *pair {
	val := i
	if i%2 == 0 {
		val = f(i)
	}
	return &pair{
		x: i,
		y: val,
	}
}

func f(x int) int {
	return x*x + x
}
