package main

import (
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/spf13/cobra"
	cliUtil "github.com/yoozoo/protocli/util"
)

const (
	langFlag = "lang"
)

type genFlagData struct {
	langValue string
}

type codeGenerator interface {
	Gen(topic string) (map[string]string, error)
}

var genFlagValue genFlagData
var outputMap = make(map[string]codeGenerator)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen-simple <topic> <output dir>",
	Short: "generate code with simple data",
	Args:  cobra.ExactArgs(2),
	Run:   generateCode,
}

func generateCode(cmd *cobra.Command, args []string) {
	generator := getCodeGenerator(genFlagValue.langValue)
	output, err := generator.Gen(args[0])
	if err != nil {
		panic(err)
	}
	err = writeFiles(args[1], output)
	if err != nil {
		panic(err)
	}
}

func getCodeGenerator(name string) codeGenerator {
	if gen, ok := outputMap[name]; ok {
		return gen
	}

	err := fmt.Errorf("Output plugin not found for %s\nsupported options: %v",
		name, reflect.ValueOf(outputMap).MapKeys())
	cliUtil.Die(err)
	return nil
}

func writeFiles(dir string, output map[string]string) error {
	outputDir := path.Dir(dir)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	for fileName, content := range output {
		f, err := os.Create(outputDir + string(os.PathSeparator) + fileName)
		if err != nil {
			return err
		}
		_, err = f.WriteString(content)
		if err != nil {
			return err
		}
		fmt.Printf("%s generated \n", fileName)
		f.Close()
	}
	return nil
}

func init() {
	genCmd.Flags().StringVar(&genFlagValue.langValue, langFlag, "php", "language of the generated code, default is php.")
}
