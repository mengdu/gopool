# Gopool

A lightweight goroutine pool for Golang.

[Reference Implementation](https://medium.com/free-code-camp/million-websockets-and-go-cc58418460bb)

```sh
go get github.com/mengdu/gopool
```

### Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/mengdu/gopool"
)

func main() {
	pool := gopool.New(3)
	defer pool.Release()
	for i := 0; i < 10; i++ {
		index := i
		// fmt.Printf("task-%d\n", index)
		pool.Schedule(func() {
			fmt.Println("do", index)
			time.Sleep(time.Second * 2) // Simulate time-consuming tasks
		})
	}
	pool.Wait() // Wait for all tasks to complete
	fmt.Println("all done")
}
```
