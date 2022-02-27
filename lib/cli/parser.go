package cli

import (
	"encoding/json"
	"fmt"
	"github.com/hellflame/argparse"
	"github.com/myOmikron/q-plugins/lib/state"
	"io/ioutil"
	"os"
	"strings"
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

func NewCommandLineInterface(
	pluginName string,
	pluginDescription string,
	pluginVersion string,
	epilog string,
) *commandLineInterface {

	cli := commandLineInterface{
		Parser: *argparse.NewParser(pluginName, pluginDescription, &argparse.ParserConfig{
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

	cli.DebugFlag = cli.Parser.Flag("d", "debug", &argparse.Option{
		Group: "general options",
		Help:  "Specify to enable debug output of the plugin",
	})

	return &cli
}

func (cli *commandLineInterface) ParseArgs() {
	if err := cli.Parser.Parse(nil); err != nil {

		switch {
		case *cli.versionFlag:
			fmt.Printf("Version: %v", cli.pluginVersion)
			os.Exit(int(state.OK))

		case *cli.generateDescriptionFlag:
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
				fmt.Println(err.Error())
				os.Exit(int(state.UNKNOWN))
			} else {
				description = string(j)
			}

			if executable, err := os.Executable(); err != nil {
				fmt.Println(err.Error())
				os.Exit(int(state.UNKNOWN))
			} else {
				if err := ioutil.WriteFile(executable+".json", []byte(description), 0664); err != nil {
					fmt.Println(err.Error())
					os.Exit(int(state.UNKNOWN))
				}
				fmt.Printf("Description was saved to %s\n", executable+".json")
				os.Exit(int(state.OK))
			}
		default:
			fmt.Println(err.Error())
			os.Exit(int(state.UNKNOWN))
		}
	}
}
