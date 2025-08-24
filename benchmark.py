import argparse
import os
import time
import json

def track_logs(pid):
    log_file = "data.json"   # same file your Go code writes
    print(f"[+] Tracking logs for PID: {pid}... Press Ctrl+C to stop.")

    last_count = 0
    while True:
        try:
            with open(log_file, "r") as f:
                data = json.load(f)

            logs = data.get("logs", {})
            current_count = len(logs)
            print(current_count)

            # Calculate logs per second over a 5 sec interval
            new_logs = current_count - last_count
            log_rate = new_logs / 1.0
            print(f"Log rate: {log_rate:.2f} logs/sec")
            last_count = current_count
            time.sleep(1)

        except KeyboardInterrupt:
            print("\n[+] Stopping benchmark.")
            break
        except FileNotFoundError:
            print(f"❌ Log file {log_file} not found")
            break
        except json.JSONDecodeError:
            # File may be mid-write by Go process
            print("⚠️ Could not parse JSON, retrying...")
            time.sleep(5)

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--pid", type=int, help="PID of process to benchmark", required=True)
    args = parser.parse_args()

    pid = args.pid
    print(f"[+] Using PID: {pid}")
    track_logs(pid)
