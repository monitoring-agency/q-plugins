package main

import (
	"errors"
	"github.com/hellflame/argparse"
	"github.com/myOmikron/q-plugins/lib/cli"
	"regexp"
)

var description = `Linux SNMP Plugin`

func loadValidator(arg string) (err error) {
	if match, _ := regexp.Match("^([0-9]+|[0-9]+\\.[0-9]*),([0-9]+|[0-9]+\\.[0-9]*),([0-9]+|[0-9]+\\.[0-9]*)$", []byte(arg)); !match {
		err = errors.New("invalid syntax for load, must be: \"{load1},{load5},{load15}\"")
	}
	return
}

func main() {
	parser := cli.NewCommandLineInterface("Linux SNMP Plugin", description, "0.1.0", "")

	loadParser := *parser.AddSubCommand("load", "Checks the load of a target", "")
	loadSnmpOptions := loadParser.AddSnmpArguments()
	loadHostname := loadParser.Parser.String("H", "hostname", &argparse.Option{
		Required: true,
		Group:    "plugin options",
		Help:     "The hostname to query.",
	})
	loadWarning := loadParser.Parser.String("", "warning", &argparse.Option{
		Group:    "plugin options",
		Help:     "Load values that determine if warning should be set. Format: load1,load5,load15",
		Validate: loadValidator,
	})
	loadCritical := loadParser.Parser.String("", "critical", &argparse.Option{
		Group:    "plugin options",
		Help:     "Load values that determine if critical should be set. Format: load1,load5,load15",
		Validate: loadValidator,
	})

	uptimeParser := *parser.AddSubCommand("uptime", "Checks the uptime of the target", "")
	uptimeSnmpOptions := uptimeParser.AddSnmpArguments()
	uptimeHostname := uptimeParser.Parser.String("H", "hostname", &argparse.Option{
		Required: true,
		Group:    "plugin options",
		Help:     "The hostname to query.",
	})

	parser.ParseArgs()

	switch {
	case loadParser.Parser.Invoked:
		cli.CheckSnmpOptions(loadSnmpOptions)
		checkLoad(loadHostname, loadWarning, loadCritical, loadSnmpOptions)
	case uptimeParser.Parser.Invoked:
		cli.CheckSnmpOptions(uptimeSnmpOptions)
		checkUptime(uptimeHostname, uptimeSnmpOptions)
	}
}
