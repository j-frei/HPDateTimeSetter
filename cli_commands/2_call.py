import datetime

def get_current_date():
    now = datetime.datetime.now()
    # Last sec entry is probably 0
    year = int(now.strftime("%y"))
    month = int(now.strftime("%m"))
    day = int(now.strftime("%d"))
    # Tricky: Tuesday was '3' and value range appears to be between 1-7
    day_of_week = ((int(datetime.datetime.now().strftime("%w")) + 7) % 7 + 1)
    hour = int(now.strftime("%H"))
    min = int(now.strftime("%M"))
    sec = int(now.strftime("%S"))
    date_array = [year, month, day, day_of_week, hour, min, sec]

    return "0x" + bytearray(date_array).hex().upper()

if __name__ == "__main__":
    import sys
    if len(sys.argv) > 1:
        host = sys.argv[1]
    else:
        host = "192.168.178.15"
        
    oid = "1.3.6.1.4.1.11.2.3.9.4.2.1.1.2.17.0"
    print(f"snmpset -v1 -c internal {host} {oid} x {get_current_date()}")
