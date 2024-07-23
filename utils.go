package main

import (
	"log"
	"os"
	"io"
	"io/fs"
	"runtime"
	"golang.org/x/sys/windows"
	"github.com/alessio/shellescape"
	"syscall"
	"strings"
)

// Copy util
func copyFile(src string, dst string) error {
	// from: https://github.com/golang/go/issues/56172#issue-1405655224
	if dst == src {
			return fs.ErrInvalid
	}

	srcF, err := os.Open(src)
	if err != nil {
			return err
	}
	defer srcF.Close()

	info, err := srcF.Stat()
	if err != nil {
			return err
	}

	dstF, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
			return err
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
			return err
	}
	return nil
}

// Windows-specific stuff
func askForAdministrativePrivilegesOnWindows() {
	// Ask for administrative privileges
	// from: https://gist.github.com/jerblack/d0eb182cc5a1c1d92d92a4c4fcc416c6
	if runtime.GOOS == "windows" {
		// Check if we are already running with elevated privileges

		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		if err == nil {
			return
		}

		// If we are not running with elevated privileges, re-run the program with elevated privileges
		verb := "runas"
		exe, err := os.Executable()
		if err != nil {
			log.Fatalf("Failed to get executable path: %v", err)
		}
		cwd, _ := os.Getwd()
		args_qt := make([]string, len(os.Args)-1)

		for i, arg := range os.Args[1:] {
			args_qt[i] = shellescape.Quote(arg)
		}
		args := strings.Join(args_qt, " ")

		verbPtr, _ := syscall.UTF16PtrFromString(verb)
		exePtr, _ := syscall.UTF16PtrFromString(exe)
		cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
		argPtr, _ := syscall.UTF16PtrFromString(args)

		var showCmd int32 = 1 //SW_NORMAL

		err = windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
		if err != nil {
			log.Fatalf("Failed to escalate privileges: %v", err)
		}
		os.Exit(0)
	}
}