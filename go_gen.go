package main

import (
	"bytes"
	"text/template"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/yoozoo/protocli/generator/data"
	"github.com/yoozoo/protomq/mqcommon"
)

type goGen struct {
	clientTpl *template.Template
	topic     string
}

func (g *goGen) Init(request *plugin.CodeGeneratorRequest) {
	g.clientTpl = util.getTpl("/template/go/client.gogo")
}

func (g *goGen) genClient(packageName string, msg *data.ProtoMessage) string {
	buf := bytes.NewBufferString("")
	data := map[string]interface{}{
		"Package": packageName,
		"StrongType": !(len(msg.Proto.Field) == 1 &&
			msg.Proto.Field[0].GetType().String() == "TYPE_STRING"),
		"Name":  msg.Proto.GetName(),
		"Topic": g.topic,
	}

	if data["StrongType"].(bool) {
		data["QueueType"] = "*" + msg.Proto.GetName()
	} else {
		data["QueueType"] = "string"
	}

	err := g.clientTpl.Execute(buf, data)
	if err != nil {
		util.Die(err)
	}

	return util.formatBuffer(buf)
}

func (g *goGen) Gen(applicationName string, packageName string, services []*data.ServiceData, messages []*data.MessageData, enums []*data.EnumData, options data.OptionMap) (result map[string]string, err error) {
	result = make(map[string]string)
	for _, msg := range messages {
		m, _ := data.GetMessageProtoAndFile(msg.Name)
		opt := m.Proto.GetOptions()
		val, err := proto.GetExtension(opt, mqcommon.E_Topic)

		// only care about msg with mqcommon.topic option
		if err != nil {
			continue
		}

		v := val.(*string)
		g.topic = *v
		content := g.genClient(applicationName, m)

		result[m.Proto.GetName()+".client.go"] = content
		println("gen", *v)
	}
	return
}

func init() {
	data.RegisterCodeGenerator("go", &goGen{})
}
