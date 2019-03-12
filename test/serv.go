package main

import (
	"context"
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/spiral/roadrunner"
)

func main() {
	ctx := context.Background()

	srv := roadrunner.NewServer(
		&roadrunner.ServerConfig{
			Command: "php phpclient/client.php echo pipes",
			Relay:   "pipes",
			Pool: &roadrunner.Config{
				NumWorkers:      1,
				AllocateTimeout: time.Second,
				DestroyTimeout:  time.Second,
			},
		})
	defer srv.Stop()

	srv.Start()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		GroupID:  "test",
		Topic:    "msg",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}

		res, err := srv.Exec(&roadrunner.Payload{Body: m.Value})
		if err != nil {
			panic(err)
		}
		fmt.Println(string(res.Body))

		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		r.CommitMessages(ctx, m)
	}
}
