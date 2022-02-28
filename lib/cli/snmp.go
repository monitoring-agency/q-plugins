package cli

import (
	"github.com/hellflame/argparse"
	"github.com/myOmikron/q-plugins/lib/validator"
)

type snmpOptions struct {
	SnmpVersion        *string
	SnmpPort           *int
	SnmpProtocol       *string
	SnmpCommunity      *string
	SnmpTimeout        *float64
	SnmpSecurityLevel  *string
	SnmpAuthProtocol   *string
	SnmpAuthPassphrase *string
	SnmpPrivProtocol   *string
	SnmpPrivPassphrase *string
}

func (cli *commandLineInterface) AddSnmpArguments() {
	snmpVersion := cli.Parser.String("", "snmp-version", &argparse.Option{
		Required: true,
		Group:    "snmp options",
		Help:     "SNMP Version. One of 1, 2c, 3",
		Choices:  []interface{}{"1", "2c", "3"},
	})

	snmpCommunity := cli.Parser.String("", "snmp-community", &argparse.Option{
		Group: "snmp options",
		Help:  "Community of SNMP. Only applies for SNMP v2c",
	})

	snmpSecurityLevel := cli.Parser.String("", "snmp-security-level", &argparse.Option{
		Group:   "snmp options",
		Help:    "Security level of SNMP messages. One of noAuthNoPriv, authNoPriv, authPriv. Only applies for SNMP v3",
		Choices: []interface{}{"noAuthNoPriv", "authNoPriv", "authPriv"},
	})
	snmpAuthenticationProtocol := cli.Parser.String("", "snmp-auth-protocol", &argparse.Option{
		Default: "sha",
		Group:   "snmp options",
		Help:    "Protocol for authentication. One of sha, md5. Only applies for SNMP v3. Defaults to sha",
		Choices: []interface{}{"sha", "md5"},
	})
	snmpAuthenticationPassphrase := cli.Parser.String("", "snmp-auth-pass", &argparse.Option{
		Group: "snmp options",
		Help:  "Passphrase for authentication. Only applies for SNMP v3",
	})
	snmpPrivProtocol := cli.Parser.String("", "snmp-priv-protocol", &argparse.Option{
		Default: "aes",
		Group:   "snmp options",
		Help:    "Protocol for privacy. One of aes, des. Only applies for SNMP v3. Defaults to aes",
		Choices: []interface{}{"aes", "des"},
	})
	snmpPrivPassphrase := cli.Parser.String("", "snmp-priv-pass", &argparse.Option{
		Group: "snmp options",
		Help:  "Passphrase for privacy. Only applies for SNMP v3",
	})

	snmpTimeout := cli.Parser.Float("", "snmp-timeout", &argparse.Option{
		Default:  "2.0",
		Group:    "snmp options",
		Help:     "Timeout of the SNMP query in seconds. Defaults to 2.0",
		Validate: validator.PositiveFloatValidator,
	})
	snmpPort := cli.Parser.Int("", "snmp-port", &argparse.Option{
		Default:  "161",
		Group:    "snmp options",
		Help:     "Port of the SNMP daemon. Defaults to 161",
		Validate: validator.PositiveIntegerValidator,
	})
	snmpProtocol := cli.Parser.String("", "snmp-protocol", &argparse.Option{
		Default: "udp",
		Group:   "snmp options",
		Help:    "Protocol used for transport. One of udp, tcp. Defaults to udp",
		Choices: []interface{}{"udp", "tcp"},
	})
	cli.SnmpOptions = &snmpOptions{
		SnmpVersion:        snmpVersion,
		SnmpPort:           snmpPort,
		SnmpCommunity:      snmpCommunity,
		SnmpTimeout:        snmpTimeout,
		SnmpProtocol:       snmpProtocol,
		SnmpSecurityLevel:  snmpSecurityLevel,
		SnmpAuthProtocol:   snmpAuthenticationProtocol,
		SnmpAuthPassphrase: snmpAuthenticationPassphrase,
		SnmpPrivProtocol:   snmpPrivProtocol,
		SnmpPrivPassphrase: snmpPrivPassphrase,
	}
}
