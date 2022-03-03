package main

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"github.com/myOmikron/q-plugins/lib/cli"
	"github.com/myOmikron/q-plugins/lib/formatter"
	"github.com/myOmikron/q-plugins/lib/protocols"
	"github.com/myOmikron/q-plugins/lib/state"
)

func checkSwap(hostname *string, warningPrct *float64, criticalPrct *float64, snmpOptions *cli.SnmpOptions) {
	oids := []string{".1.3.6.1.4.1.2021.4.3.0", ".1.3.6.1.4.1.2021.4.4.0"}
	client := protocols.GetSnmpClient(hostname, snmpOptions)
	if err := client.Connect(); err != nil {
		formatter.Error(err)
	}
	defer client.Conn.Close()

	var total, free uint64

	if packet, err := client.Get(oids); err != nil {
		formatter.Error(err)
	} else {
		for _, v := range packet.Variables {
			switch v.Name {
			case ".1.3.6.1.4.1.2021.4.3.0":
				total = gosnmp.ToBigInt(v.Value).Uint64() * 1000
			case ".1.3.6.1.4.1.2021.4.4.0":
				free = gosnmp.ToBigInt(v.Value).Uint64() * 1000
			}
		}
	}
	usage := int(total - free)
	prct := float64(usage) / float64(total) * 100

	totalData := &formatter.IntDataPoint{
		Key:   "swap-total",
		Value: int(total),
		Unit:  formatter.BYTES,
	}
	usageData := &formatter.IntDataPoint{
		Key:   "swap-used",
		Value: usage,
		Unit:  formatter.BYTES,
	}

	if *criticalPrct < prct {
		formatter.FormatOutputQ(
			fmt.Sprintf("Critical swap usage: %.2f%%: %s / %s", prct, formatter.FormatBytes(int64(usage)), formatter.FormatBytes(int64(total))),
			state.CRITICAL,
			totalData,
			usageData,
		)
	} else if *warningPrct < prct {
		formatter.FormatOutputQ(
			fmt.Sprintf("Warning swap usage: %.2f%%: %s / %s", prct, formatter.FormatBytes(int64(usage)), formatter.FormatBytes(int64(total))),
			state.WARN,
			totalData,
			usageData,
		)
	}

	formatter.FormatOutputQ(
		fmt.Sprintf("Swap usage OK: %.2f%%: %s / %s", prct, formatter.FormatBytes(int64(usage)), formatter.FormatBytes(int64(total))),
		state.OK,
		totalData,
		usageData,
	)
}
