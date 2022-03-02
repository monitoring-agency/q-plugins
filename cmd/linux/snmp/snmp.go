package main

import (
	"errors"
	"fmt"
	"github.com/hellflame/argparse"
	"github.com/myOmikron/q-plugins/lib/cli"
	"github.com/myOmikron/q-plugins/lib/formatter"
	"github.com/myOmikron/q-plugins/lib/protocols"
	"github.com/myOmikron/q-plugins/lib/state"
	"math"
	"regexp"
	"strconv"
)

func checkLoad(hostname *string, loadWarning *string, loadCritical *string, snmpOptions *cli.SnmpOptions) {
	oids := []string{".1.3.6.1.4.1.2021.10.1.3.1", ".1.3.6.1.4.1.2021.10.1.3.2", ".1.3.6.1.4.1.2021.10.1.3.3"}
	var load1, load5, load15, warnLoad1, warnLoad5, warnLoad15, critLoad1, critLoad5, critLoad15 float64
	client := protocols.GetSnmpClient(hostname, snmpOptions)

	if regex, err := regexp.Compile("^([0-9]+|[0-9]+\\.[0-9]*),([0-9]+|[0-9]+\\.[0-9]*),([0-9]+|[0-9]+\\.[0-9]*)$"); err != nil {
		formatter.Error(err)
	} else {
		parseLoad := func(r regexp.Regexp, str string) (l1, l5, l15 float64) {
			l1, l5, l15 = math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
			l := r.FindAllString(str, -1)
			for idx, v := range l {
				switch idx {
				case 0:
					l1, _ = strconv.ParseFloat(v, 64)
				case 1:
					l5, _ = strconv.ParseFloat(v, 64)
				case 2:
					l15, _ = strconv.ParseFloat(v, 64)
				}
			}
			return
		}
		warnLoad1, warnLoad5, warnLoad15 = parseLoad(*regex, *loadWarning)
		critLoad1, critLoad5, critLoad15 = parseLoad(*regex, *loadCritical)
	}

	if err := client.Connect(); err != nil {
		formatter.Error(err)
	}
	defer client.Conn.Close()

	if packet, err := client.Get(oids); err != nil {
		formatter.Error(err)
	} else {
		for idx, variable := range packet.Variables {
			if v, err := strconv.ParseFloat(string(variable.Value.([]byte)), 64); err != nil {
				switch idx {
				case 0:
					load1 = -1
				case 1:
					load5 = -1
				case 2:
					load15 = -1
				}
			} else {
				switch idx {
				case 0:
					load1 = v
				case 1:
					load5 = v
				case 2:
					load15 = v
				}
			}
		}
	}

	dataLoad1 := &formatter.Float64DataPoint{
		Key:   "load1",
		Value: load1,
	}
	dataLoad5 := &formatter.Float64DataPoint{
		Key:   "load5",
		Value: load5,
	}
	dataLoad15 := &formatter.Float64DataPoint{
		Key:   "load15",
		Value: load15,
	}

	if load1 > critLoad1 || load5 > critLoad5 || load15 > critLoad15 {
		formatter.FormatOutputQ(
			fmt.Sprintf("Critical load: 1min: %.2f, 5min: %.2f, 15min: %.2f", load1, load5, load15),
			state.CRITICAL,
			dataLoad1, dataLoad5, dataLoad15,
		)
	} else if load1 > warnLoad1 || load5 > warnLoad5 || load15 > warnLoad15 {
		formatter.FormatOutputQ(
			fmt.Sprintf("Warning load: 1min: %.2f, 5min: %.2f, 15min: %.2f", load1, load5, load15),
			state.WARN,
			dataLoad1, dataLoad5, dataLoad15,
		)
	}

	formatter.FormatOutputQ(
		fmt.Sprintf("Load OK: 1min: %.2f, 5min: %.2f, 15min: %.2f", load1, load5, load15),
		state.OK,
		dataLoad1, dataLoad5, dataLoad15,
	)
}

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

	parser.ParseArgs()

	switch {
	case loadParser.Parser.Invoked:
		cli.CheckSnmpOptions(loadSnmpOptions)
		checkLoad(loadHostname, loadWarning, loadCritical, loadSnmpOptions)
	}
}
