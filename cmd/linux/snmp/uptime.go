package main

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"github.com/myOmikron/q-plugins/lib/cli"
	"github.com/myOmikron/q-plugins/lib/formatter"
	"github.com/myOmikron/q-plugins/lib/protocols"
	"github.com/myOmikron/q-plugins/lib/state"
	"time"
)

func checkUptime(hostname *string, snmpOptions *cli.SnmpOptions) {
	oid := []string{".1.3.6.1.2.1.25.1.1.0"}
	client := protocols.GetSnmpClient(hostname, snmpOptions)
	if err := client.Connect(); err != nil {
		formatter.Error(err)
	}
	defer client.Conn.Close()

	if packet, err := client.Get(oid); err != nil {
		formatter.Error(err)
	} else {
		for _, v := range packet.Variables {
			switch v.Type {
			case gosnmp.TimeTicks:
				ticks := gosnmp.ToBigInt(v.Value).Uint64()
				uptime := time.Millisecond * time.Duration(ticks) * 10
				formatter.FormatOutputQ(fmt.Sprintf(
					"Uptime: %s", formatter.FormatTimeDuration(&uptime)),
					state.OK,
					&formatter.TimeDataPoint{
						Key:   "uptime",
						Value: uptime,
					},
				)
			}
		}
	}
}
