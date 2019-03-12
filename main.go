package main

//go:generate protoc protomq.proto --go_out=.\mqcommon

import (
	"github.com/yoozoo/protocli"
)

func main() {
	protocli.Init("protomq", "0.0.1")
	protocli.KeepDefaultLangOut = true
	protocli.Run()
}
