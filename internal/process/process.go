package process

import (
    "encoding/json"
    "sync"
    "time"
)

type ProcessInfo struct {
    PID     int    `json:"pid"`
    Status  int    `json:"status"`
    CWD     string `json:"cwd"`
    EXE     string `json:"exe"`
    Cmdline string `json:"cmdline"`
    MEM     string `json:"mem"`
    IO      string `json:"io"`
    SYSCALL string `json:"syscall"`
    CPU     string `json:"cpu"`
    FD      int    `json:"fd"`
    FDmp    string `json:"fdmp"`
}

type Process struct {
    PID  int                          `json:"pid"`
    Logs map[time.Time]ProcessInfo     `json:"logs"`
    Mu   sync.Mutex                   `json:"-"`
}

func NewProcess(pid int) *Process {
	return &Process{
		PID:  pid,
		Logs: make(map[time.Time]ProcessInfo),
	}
}

func (p *Process) ToJSON() ([]byte, error) {
    p.Mu.Lock()
    defer p.Mu.Unlock()
    return json.MarshalIndent(p,"", "    ")
}