package main

import (
	"time"
	"log"
	"flag"
	"runtime"
	"golang.org/x/sys/windows/svc"
)

// Availabilty loop
func waitForAvailability(host string, port uint16, snmp_timeout int, snmp_retries int, ping_timeout int, ping_interval int, finished chan<- bool) {
	// Wait for the host to be reachable
	log.Printf("Waiting for host %s to be reachable...", host)
	for {
		if checkIfHostIsReachable(host, ping_timeout) {
			log.Printf("Host %s is reachable", host)
			break
		}
		time.Sleep(time.Duration(ping_interval) * time.Second)
	}

	// Wait a bit more...
	time.Sleep(5 * time.Second)

	// Try to set the DateTime via SNMP
	setDateTimeViaSNMP(host, port, snmp_timeout, snmp_retries)

	// Indicate quit
	finished <- true
}


func main() {
	// Parse CLI parameters
	host := flag.String("host", "", "SNMP agent host")
	port := flag.Uint("port", 161, "SNMP agent port")
	snmp_timeout := flag.Int("snmp_timeout", 5, "Timeout in seconds for SNMP requests")
	snmp_retries := flag.Int("snmp_retries", 5, "Number of retries for SNMP requests")
	ping_timeout := flag.Int("ping_timeout", 5, "Timeout in seconds for reachability check")
	ping_loop := flag.Int("ping_interval", 60*5, "Interval in seconds for reachability check")
	mode := flag.String("mode", "", "Possible choices: install, uninstall, standalone")

	flag.Parse()

	if *host == "" {
		log.Fatalf("No host specified")
	}

	// Check if mode is in the list
	if (*mode == "") || (*mode == "standalone") {
		inService, err := false, error(nil)

		// Check if we are running as service
		if runtime.GOOS == "windows" {
			inService, err = svc.IsWindowsService()
			if err != nil {
				log.Fatalf("failed to determine if we are running in service: %v", err)
			}
		}

		if inService {
			// Run as service
			runService(*host, uint16(*port), *snmp_timeout, *snmp_retries, *ping_timeout, *ping_loop)
		} else {
			// Run in foreground
			quit_signal := make(chan bool, 1)
			waitForAvailability(*host, uint16(*port), *snmp_timeout, *snmp_retries, *ping_timeout, *ping_loop, quit_signal)
		}
	} else if *mode == "install" {
		if runtime.GOOS != "windows" {
			log.Fatalf("Service installation is only supported on Windows")
		}
		askForAdministrativePrivilegesOnWindows()

		setupService(*host, uint16(*port), *snmp_timeout, *snmp_retries, *ping_timeout, *ping_loop)
		startService()
		log.Printf("Service installed and started.")
		time.Sleep(5 * time.Second)
	} else if *mode == "uninstall" {
		if runtime.GOOS != "windows" {
			log.Fatalf("Service installation is only supported on Windows")
		}
		askForAdministrativePrivilegesOnWindows()

		removeService()
		log.Printf("Service uninstalled. Remove the directory in the Program Files folder manually.")
		time.Sleep(5 * time.Second)
	} else {
		log.Fatalf("Invalid mode specified")
		time.Sleep(5 * time.Second)
	}
}
