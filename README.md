# SetDateTime Service for older HP printers using SNMP
This repository concerns the implementation of a Windows background service to set the date and time on a HP network printer using the SNMP protocol.

Since the old, proprietary HP software solution (including `HP Solution Center`) is EOL and is partly broken, especially for more modern Windows versions, and heavily relies on Adobe Flash (EOL as well), this application replaces the time synchronization part.

The setting of the correct time is crucial in certain contexts, e.g. for a valid Fax transmission log with respect to legal aspects.

The old HP printer (`HP Photosmart 2610`) does not have an actual RTC and loses the date/time information after power loss (=pulling the plug).

## Key Contribution
This repository implements a small Go binary that 
- can run itself as a background service on Windows,
- periodically polls / waits for the availability of the network printer (default: every 5 minutes)
- sets the current date and time via SNMP as soon as the printer is available
The service terminates itself after setting the date and time.
The service is started automatically on startup.

## How to Use
### Interactive / Debug / Test Mode
1. Get the x64 executable file from the [release section](https://github.com/j-frei/HPDateTimeSetter/releases/latest) or compile it by yourself.

Compile the go binary for Windows (e.g. on Linux) using:
```bash
GOOS=windows GOARCH=amd64 go build \
    -ldflags "-w -s" \
    -o HPDateTimeSetter.exe \
    snmp_setdatetime.go
```

2. Then, run the program on Windows, using the IP or hostname of the printer (e.g. `192.168.178.15`):
```bash
.\HPDateTimeSetter.exe -host 192.168.178.15 -mode standalone
```
You can stop the binary using CTRL+C.

### Install
1. Get the x64 executable file from the [release section](https://github.com/j-frei/HPDateTimeSetter/releases/latest) or compile it by yourself.

Compile the go binary for Windows (e.g. on Linux) using:
```bash
GOOS=windows GOARCH=amd64 go build \
    -ldflags "-w -s" \
    -o HPDateTimeSetter.exe
```

2. Then, run the program on Windows, using the IP or hostname of the printer (e.g. `192.168.178.15`):
```bash
.\HPDateTimeSetter.exe -host 192.168.178.15 -mode install
```
Other CLI parameters are specified in the file `main.go`.

The service should be set up and running in the background.

### Uninstall
Run the same binary (or use the binary at `C:\Program Files\HPDateTimeService\HPDateTimeSetter.exe` with the following command in the command line (as administrator):
```bash
.\HPDateTimeSetter.exe -host 192.168.178.15 -mode uninstall
taskkill /F /IM HPDateTimeSetter.exe
rd /s /q "C:\Program Files\HPDateTimeService"
```

## About Setting DateTime via SNMP and CLI
Two different SNMPSet commands are attempted by the HP software.
Hereby, (at least) two different OIDs are used.

The calls can be generated using the following scripts at:
- `cli_commands/1_call.py` (Does not work for HP Photosmart 2610)
- `cli_commands/2_call.py` (Does work for HP Photosmart 2610)

The script only serve for debugging and educational purposes.

The Go binary avoids the installation of external SNMP tools on Windows and includes the necessary SNMP handling using static linking.