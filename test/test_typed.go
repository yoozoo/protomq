package main

import (
	test "github.com/yoozoo/protomq/temp"
)

func main() {
	c := test.NewTypedQueueMQClient([]string{"localhost:9092"})

	msg := &test.TypedQueue{
		Data: &test.Log{
			Msg:     "my msg",
			Version: 99,
		},
	}
	c.Send(msg)
}
