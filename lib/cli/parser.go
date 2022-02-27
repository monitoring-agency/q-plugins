package cli

import (
	"encoding/json"
	"fmt"
	"github.com/hellflame/argparse"
	"io/ioutil"
	"os"
)

type commandLineInterface struct {
	Parser                  argparse.Parser
	DebugFlag               *bool
	versionFlag             *bool
	generateDescriptionFlag *bool
	pluginVersion           string
	pluginName              string
	pluginDescription       string
}

func NewCommandLineInterface(applicationName string, applicationDescription string, pluginVersion string) *commandLineInterface {
	cli := commandLineInterface{
		Parser:        *argparse.NewParser(applicationName, applicationDescription, &argparse.ParserConfig{}),
		pluginVersion: pluginVersion,
	}

	cli.DebugFlag = cli.Parser.Flag("d", "debug", &argparse.Option{
		Group: "general options",
		Help:  "Specify to enable debug output of the plugin",
	})
	cli.versionFlag = cli.Parser.Flag("V", "version", &argparse.Option{
		Group: "general options",
		Help:  "Specify to output the version of the plugin and exit",
	})
	cli.generateDescriptionFlag = cli.Parser.Flag("", "generate-description", &argparse.Option{
		Group: "general options",
		Help:  "Generate the description of this plugin and save it to",
	})

	return &cli
}

func (cli *commandLineInterface) ParseArgs() {
	if err := cli.Parser.Parse(nil); err != nil {

		switch {
		case *cli.versionFlag:
			fmt.Printf("Version: %v", cli.pluginVersion)
			os.Exit(3)

		case *cli.generateDescriptionFlag:
			var description string
			if j, err := json.Marshal(struct {
				version     string
				name        string
				description string
			}{
				version:     cli.pluginVersion,
				name:        cli.pluginName,
				description: cli.pluginDescription,
			}); err != nil {
				fmt.Println(err.Error())
				os.Exit(3)
			} else {
				description = string(j)
			}

			if executable, err := os.Executable(); err != nil {
				fmt.Println(err.Error())
				os.Exit(3)
			} else {
				if err := ioutil.WriteFile(executable+".json", []byte(description), 0664); err != nil {
					fmt.Println(err.Error())
					os.Exit(3)
				}
				fmt.Printf("Description was saved to %s\n", executable)
				os.Exit(0)
			}
		default:
			fmt.Println(err.Error())
			os.Exit(3)
		}
	}
}
