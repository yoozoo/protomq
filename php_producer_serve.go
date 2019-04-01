package main

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"strings"
	"sync"

	"github.com/spiral/goridge"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

type prodServFlag struct {
	brokers string
	port    string
}

type PHPMessage struct {
	Topic   string
	Content string
}

var (
	prodFlagValue prodServFlag
	writer        *kafka.Writer
	prodServCmd   = &cobra.Command{
		Use:   "producerd",
		Short: "producerd is a php producer server",
		Args:  cobra.ExactArgs(0),
		Run:   prodServ,
	}
)

type Sender struct {
	writerMap map[string]*kafka.Writer
	wLock     sync.Mutex
}

func (s *Sender) Send(data PHPMessage, r *string) error {
	fmt.Printf("message at topic %s\n", data.Topic)

	// get kafka writer
	writer := s.getWriter(data.Topic)

	byteData := []byte(data.Content)
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Value: byteData,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) getWriter(topic string) (writer *kafka.Writer) {
	s.wLock.Lock()
	defer s.wLock.Unlock()

	writer, found := s.writerMap[topic]
	if !found {
		writer = kafka.NewWriter(kafka.WriterConfig{
			Brokers: strings.Split(prodFlagValue.brokers, ","),
			Topic:   topic,
		})
		s.writerMap[topic] = writer
	}
	return
}

func prodServ(prodServCmd *cobra.Command, args []string) {
	port := prodFlagValue.port
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	ln, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	rpc.Register(&Sender{
		writerMap: make(map[string]*kafka.Writer),
	})

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeCodec(goridge.NewCodec(conn))
	}
}

func init() {
	prodServCmd.Flags().StringVar(&prodFlagValue.brokers, "brokers", "localhost:9092", "brokers, separated by ,")
	prodServCmd.Flags().StringVar(&prodFlagValue.port, "port", ":8080", "port number")
}
