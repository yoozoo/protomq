package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/spiral/roadrunner"
)

// Consumer represents a Sarama consumer group consumer
type phpConsumer struct {
	ready  chan bool
	script string
	srv    *roadrunner.Server
	offset map[int32]int64
	group  string
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *phpConsumer) Setup(session sarama.ConsumerGroupSession) error {
	c.srv = roadrunner.NewServer(
		&roadrunner.ServerConfig{
			Command: "php " + c.script,
			Relay:   "pipes",
			Pool: &roadrunner.Config{
				NumWorkers:      consServFlagValue.numWorkers,
				AllocateTimeout: time.Second,
				DestroyTimeout:  time.Second,
			},
		})
	err := c.srv.Start()
	if err != nil {
		panic(err)
	}

	if len(c.offset) > 0 {
		for p, o := range c.offset {
			session.ResetOffset(c.group, p, o, "")
		}
	}

	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *phpConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	c.srv.Stop()
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *phpConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// The `ConsumeClaim` itself is called within a goroutine
	for message := range claim.Messages() {
		go func(message *sarama.ConsumerMessage) {
			log.Printf("Message: value = %s, timestamp = %v, topic = %s, partition = %s", string(message.Value), message.Timestamp, message.Topic, message.Partition)
			body, _ := json.Marshal([]string{message.Topic, string(message.Value)})
			_, err := c.srv.Exec(&roadrunner.Payload{Body: body})
			if err != nil {
				panic(err)
			}

			session.MarkMessage(message, "")
		}(message)
	}

	return nil
}
