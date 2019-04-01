package main

import (
	test "github.com/yoozoo/protomq/test/expected/go/test"
)

func main() {
	c := test.NewJsonQueueMQClient([]string{"localhost:9092"})
	c.Send("bingo")
}
