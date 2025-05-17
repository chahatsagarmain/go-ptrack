package process

import (
	"sync"
	"time"
)

type ProcessInfo struct {
    PID     int
    Status  int
    CWD     string
    EXE     string
    Cmdline string
    MEM  string
    IO      string
    SYSCALL string
    CPU     string
	FD 		int
	FDmp	[]string

}

type Process struct {
    PID  int
    Logs map[time.Time]ProcessInfo
    Mu   sync.Mutex
}

func NewProcess(pid int) (*Process){
	return &Process{
		PID: pid,
        Logs: make(map[time.Time]ProcessInfo),
	}
} 
