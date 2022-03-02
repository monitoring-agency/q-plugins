package main

import (
	"github.com/go-ping/ping"
	"github.com/hellflame/argparse"
	"github.com/myOmikron/q-plugins/lib/cli"
	"github.com/myOmikron/q-plugins/lib/formatter"
	"github.com/myOmikron/q-plugins/lib/state"
	"github.com/myOmikron/q-plugins/lib/validator"
	"time"
)

var description = `This plugin can be used to check hosts via ping with UDP or raw sockets.

NOTE:
    UDP:
        Unprivileged pings must be enabled:
            sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
    Raw sockets:
        As raw sockets can only be used by root, this check also must be executed in the context of root.`

func main() {
	parser := cli.NewCommandLineInterface(
		"ICMP Checker",
		description,
		"0.1.0",
		"",
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
	privileged := parser.Parser.Flag("", "privileged", &argparse.Option{
		Group: "plugin options",
		Help:  "If used, raw sockets instead of UDP are used. Requires to run the binary as root",
	})

	warningCount := parser.Parser.Int("", "warning-count", &argparse.Option{
		Group:    "plugin options",
		Validate: validator.PositiveIntegerValidator,
		Help:     "Specifies how many packets can be lost until the warning state is set",
	})
	criticalCount := parser.Parser.Int("", "critical-count", &argparse.Option{
		Group:    "plugin options",
		Validate: validator.PositiveIntegerValidator,
		Help:     "Specifies how many packets can be lost until the critical state is set",
	})
	warningRtt := parser.Parser.Int("", "warning-rtt", &argparse.Option{
		Group:    "plugin options",
		Validate: validator.PositiveIntegerValidator,
		Help:     "Warning state will be set if avg-rtt is greater than this value, in ms",
	})
	criticalRtt := parser.Parser.Int("", "critical-rtt", &argparse.Option{
		Group:    "plugin options",
		Validate: validator.PositiveIntegerValidator,
		Help:     "Critical state will be set if the avg-rtt is greater than this value, in ms",
	})

	parser.ParseArgs()

	pinger, err := ping.NewPinger(*hostname)
	if err != nil {
		formatter.Error(err)
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
		case (*criticalCount > 0 && diff > *criticalCount) || (*criticalRtt > 0 && stats.AvgRtt > time.Millisecond*time.Duration(*criticalRtt)):
			formatter.FormatOutputQ(
				"Critical output",
				state.CRITICAL,
				avgRtt,
				minRtt,
				maxRtt,
				packetloss,
			)
		case (*warningCount > 0 && diff > *warningCount) || (*warningRtt > 0 && stats.AvgRtt > time.Millisecond*time.Duration(*warningRtt)):
			formatter.FormatOutputQ(
				"Warning output",
				state.WARN,
				avgRtt,
				minRtt,
				maxRtt,
				packetloss,
			)
		case diff == 0:
			formatter.FormatOutputQ(
				"OK",
				state.OK,
				avgRtt,
				minRtt,
				maxRtt,
				packetloss,
			)
		}
	}
	pinger.Count = *count
	pinger.SetPrivileged(*privileged)
	pinger.Timeout = time.Millisecond * time.Duration(1000*(*timeout*float64(*count)+*interval*float64(*count)))
	pinger.Interval = time.Millisecond * time.Duration(1000**interval)

	if err := pinger.Run(); err != nil {
		formatter.Error(err)
	}
}
