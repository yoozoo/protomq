package main

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/yoozoo/protocli/generator/data"
)

type phpProdGen struct {
	producerTpl *template.Template
}

func (g *phpProdGen) Init(request *plugin.CodeGeneratorRequest) {
	g.producerTpl = util.getTpl("/template/php/producer.gophp")
}

func (g *phpProdGen) genProducer(applicationName, packageName, className, topic string, msg *data.ProtoMessage) string {
	buf := bytes.NewBufferString("")
	if len(packageName) <= 0 {
		packageName = applicationName
	}
	data := map[string]interface{}{
		"StrongType": !(len(msg.Proto.Field) == 1 &&
			msg.Proto.Field[0].GetType().String() == "TYPE_STRING"),
		"PackageName": strings.Title(packageName),
		"ClassName":   className,
		"Topic":       topic,
	}

	if data["StrongType"].(bool) {
		data["QueueType"] = msg.Proto.GetName()
	}

	err := g.producerTpl.Execute(buf, data)
	if err != nil {
		util.Die(err)
	}

	return buf.String()
}

func (g *phpProdGen) Gen(applicationName string, packageName string, services []*data.ServiceData, messages []*data.MessageData, enums []*data.EnumData, options data.OptionMap) (result map[string]string, err error) {
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
		className := m.Proto.GetName() + "Producer"
		fileName := className + ".php"
		if len(packageName) > 0 {
			fileName = packageName + string(os.PathSeparator) + fileName
		}
		content := g.genProducer(applicationName, packageName, className, topic, m)

		result[fileName] = content
		println("gen", topic)
	}
	return
}

func (g *phpProdGen) GetLang() string {
	return "php"
}

func init() {
	data.RegisterCodeGenerator("phpproducer", &phpProdGen{})
}
