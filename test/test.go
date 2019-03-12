package main

import (
	"fmt"
	"time"

	"github.com/Wuvist/gophpmq/msg"

	"github.com/golang/protobuf/proto"
	"github.com/spiral/roadrunner"
)

func main() {
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

	m := &msg.Msg{}
	m.Msg = "hello"
	m.Version = 2019

	data, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}

	p := &roadrunner.Payload{Body: data, Context: []byte("bingo")}

	res, err := srv.Exec(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res.Body))
	fmt.Println(string(res.Context))
}
