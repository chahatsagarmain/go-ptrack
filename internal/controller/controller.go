package controller

import (
	"fmt"
	"time"

	"github.com/chahatsagarmain/go-ptrack/internal/process"
	"github.com/chahatsagarmain/go-ptrack/internal/ptracker"
)

func ControllerStart(pid int, _ int, p *process.Process) error {
    for {
        t := time.Now().Truncate(time.Second)
        fmt.Printf("\n=== Trace at %v ===\n", t)

        // Create a single goroutine to collect all data
        done := make(chan error, 1)
        
        go func() {
            // Initialize process info
            info := process.ProcessInfo{PID: pid}
            
            // Collect all data sequentially in one goroutine
            // Status
            if resInt, err := ptracker.GetStatus(pid); err == nil && resInt != 0 {
                info.Status = resInt
            } else if err != nil {
                done <- fmt.Errorf("status error: %v", err)
                return
            }

            // CommandLine
            if res, err := ptracker.GetCommandLine(pid); err == nil {
                info.Cmdline = res
            } else {
                done <- fmt.Errorf("cmdline error: %v", err)
                return
            }

            // CWD
            if res, err := ptracker.GetCwd(pid); err == nil {
                info.CWD = res
            } else {
                done <- fmt.Errorf("cwd error: %v", err)
                return
            }

            // EXE
            if res, err := ptracker.GetExe(pid); err == nil {
                info.EXE = res
            } else {
                done <- fmt.Errorf("exe error: %v", err)
                return
            }

            // IO
            if res, err := ptracker.GetIO(pid); err == nil {
                info.IO = res
            } else {
                done <- fmt.Errorf("io error: %v", err)
                return
            }

            // SysCall
            if res, err := ptracker.GetSysCall(pid); err == nil {
                info.SYSCALL = res
            } else {
                done <- fmt.Errorf("syscall error: %v", err)
                return
            }

            // Mem
            if res, err := ptracker.GetMem(pid); err == nil {
                info.MEM = res
            } else {
                done <- fmt.Errorf("mem error: %v", err)
                return
            }

            // FD
            if res, resInt, err := ptracker.GetFD(pid); err == nil {
                info.FD = resInt
                info.FDmp = res
            } else {
                done <- fmt.Errorf("fd error: %v", err)
                return
            }

            // Single lock operation to store all data
            p.Mu.Lock()
            p.Logs[t] = info
            p.Mu.Unlock()

            done <- nil
        }()

        // Wait for completion
        if err := <-done; err != nil {
            fmt.Printf("process monitoring error: %v\n", err)
            return err
        }

        fmt.Printf("=== Completed trace %d at %v ===\n", len(p.Logs), t)

        // Add timing control
        // time.Sleep(1 * time.Second)
    }
}