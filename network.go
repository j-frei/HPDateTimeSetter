package main

import (
	"time"
	"net"
)

// Reachability check
func checkIfHostIsReachable(host string, timeout_secs int) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "80"), time.Duration(timeout_secs)*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}