# go-ptrack

**go-ptrack** is a lightweight,Highly concurrent , extensible process tracker and logger for Linux, written in Go.  
It allows you to trace a running process by its PID, periodically collect detailed runtime information from `/proc/<pid>`, and save the trace logs as JSON for further analysis even as a time based data.

---

## Features

- ✅ Track any running Linux process by PID
- ⚡ Concurrent IO and Logging 
- 📋 Collects rich process information:
  - Command line
  - Current working directory
  - Executable path
  - Memory usage
  - IO statistics
  - System calls
  - Open file descriptors
  - Status and more
- 🕒 Logs all traces with timestamps
- 🕒 Realtime Logging and Writing logs every interval 
- 📄 Outputs logs as JSON for easy analysis
- 📁 Customizable output directory with `--path` flag
- 🧩 Simple CLI interface
- 🛠️ Graceful error handling and clean shutdown

---

## Requirements

- Go 1.18+
- Linux (uses `/proc` filesystem)

---

## Installation

Clone the repository and build:

```bash
git clone https://github.com/yourusername/go-ptrack.git
cd go-ptrack
go build ./cmd/ptrack/main.go
```
You can install it as a linux command:

```bash
sudo bash install.sh
```
Now you can directly use as a linux command:
```bash
ptrack <PID> --path (optional)
```
---

## Usage

### 🔍 Basic Tracking

Start tracing a process with a specific PID:

```bash
./main track 1234
```

This will start tracing the process with PID `1234` and log traces to the default directory `/tmp/ptrack/ptrack.json`.

### 📁 Custom Output Directory

You can specify a custom directory for trace logs using the `--path` flag make sure the file exists to be written to (trace.json):

```bash
./main track 1234 --path /home/user/mytraces/trace.json
```

The trace logs will be written as JSON to `/home/user/mytraces/trace.json`.

### 🔢 Show Version

```bash
./main version
```

---

## Output

- Traces are saved as a JSON file (default: `/tmp/ptrack/ptrack.json`).
- Each trace entry contains a timestamp and a snapshot of all collected process information.

**Example JSON snippet:**

```json
{
  "timestamp": "2025-05-19T10:00:00Z",
  "pid": 1234,
  "cmdline": ["/usr/bin/myapp", "--arg1"],
  "cwd": "/home/user",
  "exe": "/usr/bin/myapp",
  "memory": {
    "VmSize": "150000 kB",
    "VmRSS": "30000 kB"
  },
  "io": {
    "read_bytes": 102400,
    "write_bytes": 51200
  },
  "fds": [0, 1, 2, 3]
}
```

---

## How It Works

- The CLI launches a controller that **periodically (every second)** collects process info using the `/proc` filesystem in a concurrent manner.
- Periodically Writes the logs to the path using a indepedent go routine .
- All data is stored in memory and periodically written to the specified output directory as a JSON file.
- If the process exits or an error occurs, tracing stops and the logs are saved.

---

## Extending

- ➕ Add new collectors by editing `ptracker.go`.
- ⚙️ Add new CLI commands or flags by editing `main.go`.

---

## Troubleshooting

- ❗ **Permission denied**:  
  Make sure you have permission to read `/proc/<pid>` and write to the output directory.

- ❗ **Process does not exist**:  
  Ensure the PID is correct and the process is running.

- ❗ **Output directory issues**:  
  Use `ptrack` or another directory you have write access to.

---

---

## Contributing

Pull requests are welcome! Feel free to open issues or suggest features.
