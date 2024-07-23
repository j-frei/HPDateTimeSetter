package main

import (
	"log"
	"time"
	"github.com/gosnmp/gosnmp"
)

// Assemble DateTime
func currentDateTimeToByteArray() []byte {
	// Determine current datetime
	cTime := time.Now()

	year := int(cTime.Year()) % 100
	month := int(cTime.Month())
	day := int(cTime.Day())

	// Shift the weekday counter to start with Sunday (Tuesday == 3)
	weekday_sunday := int(cTime.Weekday()) % 7 + 1

	hour := int(cTime.Hour())
	minute := int(cTime.Minute())
	second := int(cTime.Second())

	// Convert to OctetString
	cTime_bytearray := []byte{
		uint8(year), uint8(month), uint8(day), uint8(weekday_sunday),
		uint8(hour), uint8(minute), uint8(second),
	}

	return cTime_bytearray
}

// SNMP request
func setDateTimeViaSNMP(host string, port uint16, timeout_secs int, retries int) {
	// Set the SNMP parameters
	params := &gosnmp.GoSNMP{
		Target:    host,
		Port:      port,
		Community: "internal",
		Version:   gosnmp.Version1,
		Timeout:   time.Duration(timeout_secs) * time.Second,
		Retries:   retries,
	}

	// Connect to the SNMP agent
	conErr := params.Connect()
	if conErr != nil {
		log.Fatalf("Connect() err: %v", conErr)
	}
	defer params.Conn.Close()

	// Create an SNMP set request
	snmpPDU := []gosnmp.SnmpPDU{{
		// See: http://oidref.com/1.3.6.1.4.1.11.2.3.9.4.2.1.1.2.17
		Name:  "1.3.6.1.4.1.11.2.3.9.4.2.1.1.2.17.0",
		Type:  gosnmp.OctetString,
		Value: currentDateTimeToByteArray(),
	}}

	// Send the SNMP set request
	setResult, setErr := params.Set(snmpPDU)
	if setErr != nil {
		log.Fatalf("SNMP set request failed: %v", setErr)
	} else {
		log.Printf("SNMP set request successful: %v", setResult)
	}
}
