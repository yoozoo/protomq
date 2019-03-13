package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/yoozoo/protocli/generator/data"
	"github.com/yoozoo/protomq/mqcommon"
)

type phpGen struct {
	consumerTpl *template.Template
}

func (g *phpGen) Init(request *plugin.CodeGeneratorRequest) {
	g.consumerTpl = util.getTpl("/template/php/consumer.gophp")
}

func (g *phpGen) genConsumer(packageName string, msg *data.ProtoMessage) string {
	buf := bytes.NewBufferString("")
	data := map[string]interface{}{
		"StrongType": !(len(msg.Proto.Field) == 1 &&
			msg.Proto.Field[0].GetType().String() == "TYPE_STRING"),
		"Name": msg.Proto.GetName(),
		"GBP":  strings.ToUpper(packageName[:1]) + packageName[1:],
	}

	if data["StrongType"].(bool) {
		data["QueueType"] = msg.Proto.GetName()
	}

	err := g.consumerTpl.Execute(buf, data)
	if err != nil {
		util.Die(err)
	}

	return buf.String()
}

func (g *phpGen) Gen(applicationName string, packageName string, services []*data.ServiceData, messages []*data.MessageData, enums []*data.EnumData, options data.OptionMap) (result map[string]string, err error) {
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
		content := g.genConsumer(applicationName, m)

		result[m.Proto.GetName()+"_comsumer.php"] = content
		println("gen", *v)
	}
	return
}

func init() {
	data.RegisterCodeGenerator("php", &phpGen{})
}
