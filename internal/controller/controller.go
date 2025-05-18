package controller

import (
	"fmt"
	"time"

	"github.com/chahatsagarmain/go-ptrack/internal/process"
	"github.com/chahatsagarmain/go-ptrack/internal/ptracker"
	"golang.org/x/sync/errgroup"
)

func ControllerStart(pid int, _ int, p *process.Process) error {
	for {
		t := time.Now().Truncate(time.Second)
		fmt.Printf("\n=== Trace at %v ===\n", t)

		p.Mu.Lock()
		p.Logs[t] = process.ProcessInfo{PID: pid}
		p.Mu.Unlock()

		g := new(errgroup.Group)

		// Status
		g.Go(func() error {
			resInt, err := ptracker.GetStatus(pid)
			if err != nil || resInt == 0 {
				return fmt.Errorf("status 0 for process: %v", err)
			}
			p.Mu.Lock()
			info := p.Logs[t]
			info.Status = resInt
			p.Logs[t] = info
			p.Mu.Unlock()
			return nil
		})

		// CommandLine
		g.Go(func() error {
			res, err := ptracker.GetCommandLine(pid)
			if err != nil {
				return err
			}
			p.Mu.Lock()
			info := p.Logs[t]
			info.Cmdline = res
			p.Logs[t] = info
			p.Mu.Unlock()
			return nil
		})

		// CWD
		g.Go(func() error {
			res, err := ptracker.GetCwd(pid)
			if err != nil {
				return err
			}
			p.Mu.Lock()
			info := p.Logs[t]
			info.CWD = res
			p.Logs[t] = info
			p.Mu.Unlock()
			return nil
		})

		// EXE
		g.Go(func() error {
			res, err := ptracker.GetExe(pid)
			if err != nil {
				return err
			}
			p.Mu.Lock()
			info := p.Logs[t]
			info.EXE = res
			p.Logs[t] = info
			p.Mu.Unlock()
			return nil
		})

		// IO
		g.Go(func() error {
			res, err := ptracker.GetIO(pid)
			if err != nil {
				return err
			}
			p.Mu.Lock()
			info := p.Logs[t]
			info.IO = res
			p.Logs[t] = info
			p.Mu.Unlock()
			return nil
		})

		// SysCall
		g.Go(func() error {
			res, err := ptracker.GetSysCall(pid)
			if err != nil {
				return err
			}
			p.Mu.Lock()
			info := p.Logs[t]
			info.SYSCALL = res
			p.Logs[t] = info
			p.Mu.Unlock()
			return nil
		})

		// Mem
		g.Go(func() error {
			res, err := ptracker.GetMem(pid)
			if err != nil {
				return err
			}
			p.Mu.Lock()
			info := p.Logs[t]
			info.MEM = res
			p.Logs[t] = info
			p.Mu.Unlock()
			return nil
		})

        g.Go(func() error {
			res, resInt , err := ptracker.GetFD(pid)
			if err != nil {
				return err
			}
			p.Mu.Lock()
			info := p.Logs[t]
            info.FD = resInt
			info.FDmp = res
			p.Logs[t] = info
			p.Mu.Unlock()
			return nil
		})

		// Wait for all goroutines and handle error
		if err := g.Wait(); err != nil {
			fmt.Printf("process monitoring error: %v\n", err)
			return err
		}

		p.Mu.Lock()
		fmt.Printf("=== Completed trace %d at %v ===\n", len(p.Logs), t)
		p.Mu.Unlock()

		time.Sleep(time.Second)
	}
}
