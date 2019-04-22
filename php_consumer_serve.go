package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
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
	offset     string
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
	consServCmd.Flags().BoolVar(&consServFlagValue.oldest, "oldest", false, "if kafka consumer consume initial offset from oldest")
	consServCmd.Flags().BoolVar(&consServFlagValue.verbose, "verbose", false, "if show log")
	consServCmd.Flags().StringVar(&consServFlagValue.offset, "offset", "", "manually set offset, format: 'partition0:offset0,partition1:offset1,partition2:offset2'")
}

func consServ(cmd *cobra.Command, args []string) {
	log.Println("Starting a new kafka consumer")

	ctx := context.Background()
	topics := strings.Split(args[0], ",")
	phpScript := args[1]

	// parse offset
	offset := make(map[int32]int64, 0)
	for _, popair := range strings.Split(consServFlagValue.offset, ",") {
		temp := strings.Split(popair, ":")
		pt, err := strconv.Atoi(temp[0])
		if err != nil {
			panic(fmt.Errorf("invalid offset: %s", consServFlagValue.offset))
		}
		os, err := strconv.Atoi(temp[1])
		if err != nil {
			panic(fmt.Errorf("invalid offset: %s", consServFlagValue.offset))
		}
		offset[int32(pt)] = int64(os)
	}
	consumer := &phpConsumer{
		script: phpScript,
		offset: offset,
		group:  consServFlagValue.group,
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
