package main

import(
	"fmt"

	"github.com/gude/queue"
)

func main() {
	q := queue.New(5)
	q.Enqueue("hello", "world")
	fmt.Println("len ", q.Len())
	key, val, ok := q.Dequeue()
	if ok {
		fmt.Println(key,val)
	}
	fmt.Println("len ", q.Len())
}
