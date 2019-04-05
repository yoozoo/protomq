package goconsumer

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/segmentio/kafka-go"
)

type consumerHandler func(*kafka.Message) error

type Consumer struct {
	topic   string
	group   string
	brokers string
	handler consumerHandler
}

func NewConsumer() *Consumer {
	consumer := &Consumer{}
	flag.StringVar(&consumer.topic, "topic", "", "topic in msg go")
	flag.StringVar(&consumer.group, "group", "go", "consumer group; default go")
	flag.StringVar(&consumer.brokers, "brokers", "localhost:9092", "brokers, separated by ,")
	flag.Parse()
	return consumer
}

func (c *Consumer) RegisterHandler(handler consumerHandler) {
	c.handler = handler
}

func (c *Consumer) SetTopic(topic string) {
	c.topic = topic
}

func (c *Consumer) Run() {
	ctx := context.Background()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(c.brokers, ","),
		GroupID:  c.group,
		Topic:    c.topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}
		go func() {
			fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			if err != nil {
				panic(err)
			}
			c.handler(&m)
			if err != nil {
				panic(err)
			}
			r.CommitMessages(ctx, m)
		}()
	}
}
