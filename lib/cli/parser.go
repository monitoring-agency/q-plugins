package cli

import (
	"fmt"
	"github.com/hellflame/argparse"
	"os"
)

type commandLineInterface struct {
	Parser      argparse.Parser
	DebugFlag   *bool
	versionFlag *bool
	version     string
}

func NewCommandLineInterface(applicationName string, applicationDescription string, pluginVersion string) *commandLineInterface {
	cli := commandLineInterface{
		Parser:  *argparse.NewParser(applicationName, applicationDescription, &argparse.ParserConfig{}),
		version: pluginVersion,
	}

	cli.DebugFlag = cli.Parser.Flag("d", "debug", &argparse.Option{
		Group: "general options",
		Help:  "Specify to enable debug output of the plugin",
	})
	cli.versionFlag = cli.Parser.Flag("V", "version", &argparse.Option{
		Group: "general options",
		Help:  "Specify to output the version of the plugin and exit",
	})
	return &cli
}

func (cli *commandLineInterface) ParseArgs() {
	if err := cli.Parser.Parse(nil); err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}

	if *cli.versionFlag {
		fmt.Printf("Version: %v", cli.version)
		os.Exit(3)
	}
}
