package main

//go:generate protoc protomq.proto --go_out=.\mqcommon
//go:generate esc -o tpl/tpl.go -modtime 0 -pkg=tpl template protomq.proto

import (
	"github.com/yoozoo/protocli"
	"github.com/yoozoo/protomq/tpl"
)

const commonProtoPath = "/protomq.proto"

func main() {
	protocli.Init("protomq", "0.0.1")
	protocli.KeepDefaultLangOut = true
	protocli.RootCmd.AddCommand(consServCmd)
	protocli.RootCmd.AddCommand(prodServCmd)
	protocli.RootCmd.AddCommand(genCmd)
	protocli.RegisterIncludeFile(commonProtoPath, tpl.FSMustString(false, commonProtoPath))
	protocli.Run()
}
