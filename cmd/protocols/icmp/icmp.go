package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/hellflame/argparse"
	"github.com/myOmikron/q-plugins/lib/cli"
	"github.com/myOmikron/q-plugins/lib/formatter"
	"github.com/myOmikron/q-plugins/lib/validator"
	"os"
	"time"
)

func main() {
	parser := cli.NewCommandLineInterface(
		"ICMP Checker",
		"This is a program to check hosts via icmp",
		"0.1.0",
	)

	hostname := parser.Parser.String("H", "hostname", &argparse.Option{
		Required: true,
		Group:    "plugin options",
		Help:     "The hostname the icmp packets are sent to. Can be IPv4, IPv6 or DNS / Hostname",
	})
	count := parser.Parser.Int("c", "count", &argparse.Option{
		Default:  "3",
		Validate: validator.PositiveIntegerValidator,
		Group:    "plugin options",
		Help:     "The count of packets to send",
	})
	interval := parser.Parser.Float("i", "interval", &argparse.Option{
		Default:  "1.0",
		Group:    "plugin options",
		Validate: validator.PositiveFloatValidator,
		Help:     "The interval between the packets, in seconds",
	})
	timeout := parser.Parser.Float("t", "timeout", &argparse.Option{
		Default:  "2.0",
		Group:    "plugin options",
		Validate: validator.PositiveFloatValidator,
		Help:     "The timeout of each packet, in seconds",
	})
	warningCount := parser.Parser.Int("", "warning-count", &argparse.Option{
		Default:  "1",
		Group:    "plugin options",
		Validate: validator.PositiveIntegerValidator,
		Help:     "Specifies how many packets can be lost until the warning state is set",
	})
	criticalCount := parser.Parser.Int("", "critical-count", &argparse.Option{
		Default:  "2",
		Group:    "plugin options",
		Validate: validator.PositiveIntegerValidator,
		Help:     "Specifies how many packets can be lost until the critical state is set",
	})

	parser.ParseArgs()

	pinger, err := ping.NewPinger(*hostname)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		diff := stats.PacketsSent - stats.PacketsRecv
		avgRtt := &formatter.TimeDataPoint{
			Key:   "avg-rtt",
			Value: stats.AvgRtt,
		}
		minRtt := &formatter.TimeDataPoint{
			Key:   "min-rtt",
			Value: stats.MinRtt,
		}
		maxRtt := &formatter.TimeDataPoint{
			Key:   "max-rtt",
			Value: stats.MaxRtt,
		}
		packetloss := &formatter.Float64DataPoint{
			Key:   "packetloss",
			Value: stats.PacketLoss,
		}

		switch {
		case diff > *criticalCount:
			formatter.FormatOutputQ(
				"Critical output",
				formatter.CRITICAL,
				avgRtt,
				minRtt,
				maxRtt,
				packetloss,
			)
		case diff > *warningCount:
			formatter.FormatOutputQ(
				"Warning output",
				formatter.WARN,
				avgRtt,
				minRtt,
				maxRtt,
				packetloss,
			)
		case diff == 0:
			formatter.FormatOutputQ(
				"OK",
				formatter.OK,
				avgRtt,
				minRtt,
				maxRtt,
				packetloss,
			)
		}
	}
	pinger.Count = *count
	pinger.Timeout = time.Millisecond * time.Duration(1000*(*timeout*float64(*count)+*interval*float64(*count)))
	pinger.Interval = time.Millisecond * time.Duration(1000**interval)

	if err := pinger.Run(); err != nil {
		fmt.Println(err.Error())
	}
}
