import datetime
# snmpget -v1 -c internal 192.168.178.15 1.3.6.1.4.1.11.2.3.9.4.2.1.1.2.8.0

def get_current_date():
    now = datetime.datetime.now()
    # Last sec entry is probably 0
    current_date_hex = (now.strftime("%Y%m%d%H%M") + "0").encode("utf-8").hex()
    return f"0xFDE8{current_date_hex.upper()}"

if __name__ == "__main__":
    import sys
    if len(sys.argv) > 1:
        host = sys.argv[1]
    else:
        host = "192.168.178.15"

    oid = "1.3.6.1.4.1.11.2.3.9.4.2.1.1.2.8.0"
    print(f"snmpset -v1 -c internal {host} {oid} x {get_current_date()}")