package protocols

import (
	"github.com/gosnmp/gosnmp"
	"github.com/myOmikron/q-plugins/lib/cli"
	"time"
)

func GetSnmpClient(hostname *string, snmpOptions *cli.SnmpOptions) *gosnmp.GoSNMP {
	var version gosnmp.SnmpVersion
	switch *snmpOptions.SnmpVersion {
	case "1":
		version = gosnmp.Version1
	case "2c":
		version = gosnmp.Version2c
	case "3":
		version = gosnmp.Version3
	}

	var authProtocol gosnmp.SnmpV3AuthProtocol
	switch *snmpOptions.SnmpAuthProtocol {
	case "sha":
		authProtocol = gosnmp.SHA
	case "md5":
		authProtocol = gosnmp.MD5
	}

	var privProtocol gosnmp.SnmpV3PrivProtocol
	switch *snmpOptions.SnmpPrivProtocol {
	case "des":
		privProtocol = gosnmp.DES
	case "aes":
		privProtocol = gosnmp.AES
	}
	securityModel := gosnmp.UserSecurityModel

	var msgFlag gosnmp.SnmpV3MsgFlags
	switch *snmpOptions.SnmpSecurityLevel {
	case "noAuthNoPriv":
		msgFlag = gosnmp.NoAuthNoPriv
	case "authNoPriv":
		msgFlag = gosnmp.AuthNoPriv
	case "authPriv":
		msgFlag = gosnmp.AuthPriv
	}

	return &gosnmp.GoSNMP{
		Target:        *hostname,
		Port:          uint16(*snmpOptions.SnmpPort),
		Transport:     *snmpOptions.SnmpProtocol,
		Community:     *snmpOptions.SnmpCommunity,
		Version:       version,
		Timeout:       time.Duration(*snmpOptions.SnmpTimeout) * time.Second,
		MsgFlags:      msgFlag,
		SecurityModel: securityModel,
		SecurityParameters: &gosnmp.UsmSecurityParameters{
			UserName:                 *snmpOptions.SnmpUser,
			AuthenticationProtocol:   authProtocol,
			PrivacyProtocol:          privProtocol,
			AuthenticationPassphrase: *snmpOptions.SnmpAuthPassphrase,
			PrivacyPassphrase:        *snmpOptions.SnmpPrivPassphrase,
		},
	}
}
