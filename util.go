package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"strconv"
	"text/template"

	"github.com/golang/protobuf/proto"
	"github.com/yoozoo/protocli/generator/data"

	"github.com/yoozoo/protomq/mqcommon"
	"github.com/yoozoo/protomq/tpl"
)

var (
	rgxSyntaxError = regexp.MustCompile(`(\d+):\d+: `)
	debugTpl       = os.Getenv("debugTpl") == "true"
)

var util = &_util{}

type _util struct{}

func (g *_util) getTpl(path string) *template.Template {
	var err error
	tpl := template.New("tpl")
	tplStr := g.LoadTpl(path)
	result, err := tpl.Parse(tplStr)
	if err != nil {
		g.Die(err)
	}
	return result
}

func (g *_util) formatBuffer(buf *bytes.Buffer) string {
	output, err := format.Source(buf.Bytes())
	if err == nil {
		return string(output)
	}

	matches := rgxSyntaxError.FindStringSubmatch(err.Error())
	if matches == nil {
		util.Die(errors.New("failed to format template"))
	}

	lineNum, _ := strconv.Atoi(matches[1])
	scanner := bufio.NewScanner(buf)
	errBuf := &bytes.Buffer{}
	line := 1
	for ; scanner.Scan(); line++ {
		if delta := line - lineNum; delta < -5 || delta > 5 {
			continue
		}

		if line == lineNum {
			errBuf.WriteString(">>>> ")
		} else {
			fmt.Fprintf(errBuf, "% 4d ", line)
		}
		errBuf.Write(scanner.Bytes())
		errBuf.WriteByte('\n')
	}

	util.Die(fmt.Errorf("failed to format template\n\n%s", errBuf.Bytes()))

	return ""
}

// Die prints error and exit
func (g *_util) Die(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

// LoadTpl is the function to load template file as string
// It loads file content from esc embed by default
// Set environment variable debugTpl to "true" to load template from disk directly
func (g *_util) LoadTpl(tplPath string) string {
	//useLocal is true, the filesystem's contents are instead used.
	return tpl.FSMustString(debugTpl, tplPath)
}

// RetriveTopics will return error if there is no topic
// or there are repeated topics
func (g *_util) RetriveTopics(messages []*data.MessageData) (map[string]string, error) {
	topics := make(map[string]bool, 0)
	result := make(map[string]string, 0)

	for _, msg := range messages {

		m, _ := data.GetMessageProtoAndFile(msg.Name)
		opt := m.Proto.GetOptions()
		val, err := proto.GetExtension(opt, mqcommon.E_Topic)

		if err != nil {
			continue
		}

		v := val.(*string)
		_, found := topics[*v]
		if found {
			return nil, fmt.Errorf("repeated topic: %s", *v)
		}

		topics[*v] = true
		result[msg.Name] = *v
	}

	if len(topics) <= 0 {
		return nil, fmt.Errorf("no topic found")
	}
	return result, nil
}
