package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

type consServFlag struct {
	group      string
	brokers    string
	oldest     bool
	verbose    bool
	numWorkers int64
}

var consServFlagValue consServFlag

var consServCmd = &cobra.Command{
	Use:   "consumerd <topcis> <comsumer.php>",
	Short: "consumerd serve a php script as mq consumer. \"topcis\" separated by \",\"",
	Args:  cobra.ExactArgs(2),
	Run:   consServ,
}

func init() {
	consServCmd.Flags().StringVar(&consServFlagValue.group, "group", "php", "consumer group; default php")
	consServCmd.Flags().StringVar(&consServFlagValue.brokers, "brokers", "localhost:9092", "brokers, separated by ,")
	consServCmd.Flags().Int64Var(&consServFlagValue.numWorkers, "workers", 5, "number of max php workers")
	consServCmd.Flags().BoolVar(&consServFlagValue.oldest, "oldest", true, "if kafka consumer consume initial offset from oldest")
	consServCmd.Flags().BoolVar(&consServFlagValue.verbose, "verbose", false, "if show log")
}

func consServ(cmd *cobra.Command, args []string) {
	log.Println("Starting a new kafka consumer")

	ctx := context.Background()
	topics := strings.Split(args[0], ",")
	phpScript := args[1]

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := &phpConsumer{
		script: phpScript,
	}

	config := sarama.NewConfig()

	if consServFlagValue.oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	version, _ := sarama.ParseKafkaVersion("0.10.2.0")
	config.Version = version

	client, err := sarama.NewConsumerGroup(
		strings.Split(consServFlagValue.brokers, ","),
		consServFlagValue.group,
		config,
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			consumer.ready = make(chan bool, 0)
			err := client.Consume(ctx, topics, consumer)
			if err != nil {
				panic(err)
			}
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm // Await a sigterm signal before safely closing the consumer
	consumer.srv.Stop()
	err = client.Close()
	if err != nil {
		panic(err)
	}
}
