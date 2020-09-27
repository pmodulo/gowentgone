package main

import "fmt"

type building struct {
	left   int
	right  int
	height int
}

func main() {
	inCh := make(chan building)
	quit := make(chan interface{})

	go compute(inCh, quit)

	inCh <- building{
		left:   1,
		right:  15,
		height: 3,
	}
	inCh <- building{
		left:   4,
		right:  11,
		height: 5,
	}
	inCh <- building{
		left:   19,
		right:  23,
		height: 4,
	}
	close(inCh)
	<-quit
}

func compute(in <-chan building, quit chan<- interface{}) {
	heights := make([]int, 100)

	for b := range in {
		for i := b.left; i < b.right; i++ {
			if heights[i] < b.height {
				heights[i] = b.height
			}
		}
	}
	prev := 0
	for i, v := range heights {
		if v != prev {
			fmt.Println(i, v)
			prev = v
		}
	}
	close(quit)
}
