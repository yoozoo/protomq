package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spiral/roadrunner"
)

type servFlag struct {
	topic   string
	group   string
	brokers string
}

var servFlagValue servFlag

var servCmd = &cobra.Command{
	Use:   "serv <comsumer.php>",
	Short: "serve a php script as mq comsumer",
	Args:  cobra.ExactArgs(1),
	Run:   serv,
}

func init() {
	servCmd.Flags().StringVar(&servFlagValue.topic, "topic", "", "topic in msg go")
	servCmd.Flags().StringVar(&servFlagValue.group, "group", "php", "comsumer group; default php")
	servCmd.Flags().StringVar(&servFlagValue.brokers, "brokers", "localhost:9092", "brokers, separated by ,")
}

func serv(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	phpScript := args[0]

	srv := roadrunner.NewServer(
		&roadrunner.ServerConfig{
			Command: "php " + phpScript + " echo pipes",
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
		Brokers:  strings.Split(servFlagValue.brokers, ","),
		GroupID:  servFlagValue.group,
		Topic:    servFlagValue.topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}

		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		_, err = srv.Exec(&roadrunner.Payload{Body: m.Value})
		if err != nil {
			panic(err)
		}

		r.CommitMessages(ctx, m)
	}
}
