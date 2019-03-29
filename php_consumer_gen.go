package main

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/yoozoo/protocli/generator/data"
)

type phpGen struct {
	consumerTpl *template.Template
}

func (g *phpGen) Init(request *plugin.CodeGeneratorRequest) {
	g.consumerTpl = util.getTpl("/template/php/consumer.gophp")
}

func (g *phpGen) genConsumer(applicationName, packageName, className string, msg *data.ProtoMessage) string {
	buf := bytes.NewBufferString("")
	if len(packageName) <= 0 {
		packageName = applicationName
	}
	data := map[string]interface{}{
		"StrongType": !(len(msg.Proto.Field) == 1 &&
			msg.Proto.Field[0].GetType().String() == "TYPE_STRING"),
		"Name":        msg.Proto.GetName(),
		"GBP":         strings.ToUpper(applicationName[:1]) + applicationName[1:],
		"PackageName": strings.Title(packageName),
		"ClassName":   className,
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

		m, _ := data.GetMessageProtoAndFile(msg.Name)
		className := m.Proto.GetName() + "_consumer"
		fileName := className + ".php"
		if len(packageName) > 0 {
			fileName = packageName + string(os.PathSeparator) + fileName
		}
		content := g.genConsumer(applicationName, packageName, className, m)

		result[fileName] = content
		println("gen", topic)
	}
	return
}

func init() {
	data.RegisterCodeGenerator("phpconsumer", &phpGen{})
}
