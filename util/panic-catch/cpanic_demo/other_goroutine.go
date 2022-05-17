package main

import (
	"fmt"
	"time"
)

func procPan() {
	fmt.Println("into a new goroutine")

	select {
	case <-time.After(10 * time.Second):
		{
			aryInts := make([]int, 3, 3)
			aryInts[9] = 3
		}
	}
}
