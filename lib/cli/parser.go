package cli

import (
	"encoding/json"
	"fmt"
	"github.com/hellflame/argparse"
	"github.com/myOmikron/q-plugins/lib/formatter"
	"github.com/myOmikron/q-plugins/lib/state"
	"io/ioutil"
	"os"
	"strings"
)

type commandLineInterface struct {
	Parser                  *argparse.Parser
	versionFlag             *bool
	generateDescriptionFlag *bool
	pluginVersion           string
	pluginName              string
	pluginDescription       string
}

func NewCommandLineInterface(
	pluginName string,
	pluginDescription string,
	pluginVersion string,
	epilog string,
) *commandLineInterface {

	cli := commandLineInterface{
		Parser: argparse.NewParser(pluginName, pluginDescription, &argparse.ParserConfig{
			EpiLog: epilog,
		}),
		pluginVersion:     pluginVersion,
		pluginDescription: pluginDescription,
		pluginName:        pluginName,
	}
	cli.versionFlag = cli.Parser.Flag("V", "version", &argparse.Option{
		Help: "Specify to output the version of the plugin and exit",
	})
	cli.generateDescriptionFlag = cli.Parser.Flag("", "generate-description", &argparse.Option{
		Help: "Generate the description of this plugin and save it to {executable}.json",
	})

	return &cli
}

func (cli *commandLineInterface) AddSubCommand(cmd string, commandDescription string, epilog string) *commandLineInterface {
	childParser := cli.Parser.AddCommand(cmd, commandDescription, &argparse.ParserConfig{
		Usage:  commandDescription,
		EpiLog: epilog,
	})

	return &commandLineInterface{
		Parser: childParser,
	}
}

func (cli *commandLineInterface) generateDescription() {
	splitDescription := strings.Split(cli.pluginDescription, "\n\n")
	shortDescription := splitDescription[0]

	var description string
	if j, err := json.Marshal(&struct {
		Version     string `json:"version"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}{
		Version:     cli.pluginVersion,
		Name:        cli.pluginName,
		Description: shortDescription,
	}); err != nil {
		formatter.Error(err)
	} else {
		description = string(j)
	}

	if executable, err := os.Executable(); err != nil {
		formatter.Error(err)
	} else {
		if err := ioutil.WriteFile(executable+".json", []byte(description), 0664); err != nil {
			formatter.Error(err)
		}
		fmt.Printf("Description was saved to %s\n", executable+".json")
		os.Exit(int(state.OK))
	}
}

func (cli *commandLineInterface) checkDefaultArguments() {
	switch {
	case *cli.versionFlag:
		fmt.Printf("Version: %v", cli.pluginVersion)
		os.Exit(int(state.OK))

	case *cli.generateDescriptionFlag:
		cli.generateDescription()
	}
}

func (cli *commandLineInterface) ParseArgs() {
	err := cli.Parser.Parse(nil)
	if err != nil {
		// This will terminate if a default arg was found
		cli.checkDefaultArguments()

		formatter.Error(err)
	}

	cli.checkDefaultArguments()
}
