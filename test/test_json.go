package main

import (
	test "github.com/yoozoo/protomq/temp"
)

func main() {
	c := test.NewJsonQueueMQClient([]string{"localhost:9092"})
	c.Send("bingo")
}
