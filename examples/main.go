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
			time.Sleep(time.Second * 2) // 模拟耗时任务
		})
	}
	pool.Wait()
	fmt.Println("all done")
}
