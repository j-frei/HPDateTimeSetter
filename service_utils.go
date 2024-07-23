package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
)

func setupService(host string, port uint16, snmp_timeout int, snmp_retries int, ping_timeout int, ping_interval int) {
	// Determine Program Files directory
	program_files := os.Getenv("ProgramFiles")
	dir_path := filepath.Join(program_files, INSTALLED_DIRECTORY_NAME)

	executable_path, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	// Create the directory on Windows
	err = os.MkdirAll(dir_path, 0755)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Copy the executable to the directory
	executable_target := filepath.Join(dir_path, SNMP_SETDATE_SERVICE_EXECUTABLE)
	err = copyFile(executable_path, executable_target)
	if err != nil {
		log.Fatalf("Failed to copy executable: %v", err)
	}

	// Register a new service
	err = addService(
		executable_target,
		[]string{
			"-host", host,
			"-port", strconv.Itoa(int(port)),
			"-snmp_timeout", strconv.Itoa(snmp_timeout),
			"-snmp_retries", strconv.Itoa(snmp_retries),
			"-ping_timeout", strconv.Itoa(ping_timeout),
			"-ping_interval", strconv.Itoa(ping_interval)},
	)

	if err != nil {
		log.Fatalf("Failed to add the service: %v", err)
	}
}

// Service add/remove
func addService(execPath string, cmdArgs []string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(SNMP_SETDATE_SERVICE_NAME)
	if err == nil {
		s.Close()
		log.Fatalf("Service %s already exists", SNMP_SETDATE_SERVICE_NAME)
	}

	s, err = m.CreateService(SNMP_SETDATE_SERVICE_NAME, execPath, mgr.Config{
		StartType: windows.SERVICE_AUTO_START,
		DisplayName: SNMP_SETDATE_SERVICE_DISPLAY_NAME,
		Description: SNMP_SETDATE_SERVICE_DESCRIPTION,
		DelayedAutoStart: true,
		}, cmdArgs...)

	if err != nil {
		log.Fatalf("CreateService() failed: %s", err)
	}
	defer s.Close()

	return nil
}

func removeService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(SNMP_SETDATE_SERVICE_NAME)
	if err != nil {
		log.Fatalf("Service %s is not installed", SNMP_SETDATE_SERVICE_NAME)
	} else {
		// Try to stop service
		s.Control(windows.SERVICE_CONTROL_STOP)
	}
	defer s.Close()

	err = s.Delete()
	if err != nil {
		log.Fatalf("Delete() failed: %s", err)
	}

	return nil
}

func startService() error {

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(SNMP_SETDATE_SERVICE_NAME)
	if err != nil {
		log.Fatalf("Service %s is not installed", SNMP_SETDATE_SERVICE_NAME)
	}
	defer s.Close()

	err = s.Start()
	if err != nil {
		return err
	}

	return nil
}