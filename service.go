package main

import (
	"time"
	"log"
	"golang.org/x/sys/windows/svc"
)

const SNMP_SETDATE_SERVICE_NAME = "HPDateTimeService"
const SNMP_SETDATE_SERVICE_DISPLAY_NAME = "HP DateTime Service"
const SNMP_SETDATE_SERVICE_DESCRIPTION = "This service sets the DateTime for a HP Printer using SNMP"
const SNMP_SETDATE_SERVICE_EXECUTABLE = "HPDateTimeSetter.exe"
const INSTALLED_DIRECTORY_NAME = "HPDateTimeService"

func runService(host string, port uint16, snmp_timeout int, snmp_retries int, ping_timeout int, ping_interval int) {
	err := svc.Run(SNMP_SETDATE_SERVICE_NAME, &snmpDateTimeService{
		host: host,
		port: port,
		snmp_timeout: snmp_timeout,
		snmp_retries: snmp_retries,
		ping_timeout: ping_timeout,
		ping_interval: ping_interval,
	})

	if err != nil {
		log.Fatalf("Failed to run service: %v", err)
	}
	return
}

type snmpDateTimeService struct{
	host string
	port uint16
	snmp_timeout int
	snmp_retries int
	ping_timeout int
	ping_interval int
}

func (m *snmpDateTimeService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptShutdown | svc.AcceptStop
	changes <- svc.Status{State: svc.StartPending}

	finished := make(chan bool, 1)
	go waitForAvailability(m.host, m.port, m.snmp_timeout, m.snmp_retries, m.ping_timeout, m.ping_interval, finished)
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case <-finished:
			break loop
		case c := <-r:
			switch c.Cmd {
				case svc.Interrogate:
					changes <- c.CurrentStatus
					// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
					time.Sleep(100 * time.Millisecond)
					changes <- c.CurrentStatus
				case svc.Shutdown | svc.Stop:
					break loop
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}
