package main

import (
	"context"
	"net"
	"net/rpc"
	"strings"

	"github.com/spiral/goridge"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

type prodServFlag struct {
	topic   string
	brokers string
	port    string
}

type rpcObj interface{}

var (
	prodFlagValue prodServFlag
	cmd           *cobra.Command
	writer        *kafka.Writer
	prodServCmd   = &cobra.Command{
		Use:   "prodServ <comsumer.php>",
		Short: "prodServe a php script as mq consumer",
		Args:  cobra.ExactArgs(1),
		Run:   prodServ,
	}
)

type Sender struct{}

func (s *Sender) Send(data string, r *string) error {
	byteData := []byte(data)
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Value: byteData,
	})
	if err != nil {
		return err
	}
	return nil
}

func prodServ(cmd *cobra.Command, args []string) {
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: strings.Split(prodFlagValue.brokers, ","),
		Topic:   prodFlagValue.topic,
	})

	port := prodFlagValue.port
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	ln, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	rpc.Register(new(Sender))

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeCodec(goridge.NewCodec(conn))
	}
}

func init() {
	cmd.Flags().StringVar(&prodFlagValue.topic, "topic", "", "topic in msg go")
	cmd.Flags().StringVar(&prodFlagValue.brokers, "brokers", "localhost:9092", "brokers, separated by ,")
	cmd.Flags().StringVar(&prodFlagValue.port, "port", ":8080", "port number")
}
