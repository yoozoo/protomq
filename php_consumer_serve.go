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

type consServFlag struct {
	topic   string
	group   string
	brokers string
}

var consServFlagValue consServFlag

var consServCmd = &cobra.Command{
	Use:   "consServ <comsumer.php>",
	Short: "consServe a php script as mq consumer",
	Args:  cobra.ExactArgs(1),
	Run:   consServ,
}

func init() {
	consServCmd.Flags().StringVar(&consServFlagValue.topic, "topic", "", "topic in msg go")
	consServCmd.Flags().StringVar(&consServFlagValue.group, "group", "php", "comsumer group; default php")
	consServCmd.Flags().StringVar(&consServFlagValue.brokers, "brokers", "localhost:9092", "brokers, separated by ,")
}

func consServ(cmd *cobra.Command, args []string) {
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

	err := srv.Start()
	if err != nil {
		panic(err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(consServFlagValue.brokers, ","),
		GroupID:  consServFlagValue.group,
		Topic:    consServFlagValue.topic,
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
