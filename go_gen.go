package main

import (
	"bytes"
	"text/template"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/yoozoo/protocli/generator/data"
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
	topicMap, err := util.RetriveTopics(messages)
	if err != nil {
		return nil, err
	}
	for _, msg := range messages {
		topic, found := topicMap[msg.Name]
		// only care about msg with mqcommon.topic option
		if !found {
			continue
		}

		g.topic = topic
		m, _ := data.GetMessageProtoAndFile(msg.Name)
		content := g.genClient(applicationName, m)

		result[m.Proto.GetName()+".client.go"] = content
		println("gen", topic)
	}
	return
}

func init() {
	data.RegisterCodeGenerator("go", &goGen{})
}
