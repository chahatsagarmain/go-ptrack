package controller

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chahatsagarmain/go-ptrack/internal/process"
	"github.com/chahatsagarmain/go-ptrack/internal/ptracker"
)


func ControllerStart(pid int , _ int , p *process.Process) (error) {

	var wg sync.WaitGroup;

	errChan := make(chan error);

	for {
		
		t := time.Now().Truncate(time.Second);
		fmt.Printf("\n=== Trace at %v ===\n", t)
		p.Mu.Lock();
		p.Logs[t] = process.ProcessInfo{
			PID: pid,
		}
		p.Mu.Unlock();
		var res string;
		var err error;

		wg.Add(1);
		go func(tn time.Time) {
			var resInt int;
			resInt , err = ptracker.GetStatus(pid);
			if err != nil || resInt == 0 {
				errChan <- err;
				log.Fatal(err);
			}
			fmt.Printf("%v\n",resInt);
			p.Mu.Lock();
			info := p.Logs[tn];
			info.Status = resInt;
			p.Logs[tn] = info;
			p.Mu.Unlock();
			wg.Done();
		}(t);

		wg.Add(1)
		go func(tn time.Time) {
			res , err = ptracker.GetCommandLine(pid)
			if err != nil {
				errChan <- err;
				log.Fatal(err);
			}	
			fmt.Printf("%v\n",res);
			p.Mu.Lock();
			info := p.Logs[tn];
			info.Cmdline = res;
			p.Logs[tn] = info;
			p.Mu.Unlock();
			wg.Done();
		}(t);

		wg.Add(1);
		go func(tn time.Time) {
			res , err = ptracker.GetCwd(pid);
			if err != nil {
				errChan <- err;
				log.Fatal(err);
			}	
			fmt.Printf("%v\n",res);
			p.Mu.Lock();
			info := p.Logs[tn];
			info.CWD = res;
			p.Logs[tn] = info;
			p.Mu.Unlock();
			wg.Done();
		}(t);

		wg.Add(1);
		go func(tn time.Time) {
			res , err = ptracker.GetExe(pid);
			if err != nil {
				errChan <- err;
				log.Fatal(err);
			}	
			fmt.Printf("%v\n",res);
			p.Mu.Lock();
			info := p.Logs[tn];
			info.EXE = res;
			p.Logs[tn] = info;
			p.Mu.Unlock();
			wg.Done();
		}(t);

		wg.Add(1);
		go func(tn time.Time) {
			res , err = ptracker.GetIO(pid);
			if err != nil {
				errChan <- err;
				log.Fatal(err);
			}	
			fmt.Printf("%v\n",res);
			p.Mu.Lock();
			info := p.Logs[tn];
			info.IO = res;
			p.Logs[tn] = info;
			p.Mu.Unlock();
			wg.Done();
		}(t);

		wg.Add(1);
		go func(tn time.Time) {
			res , err = ptracker.GetSysCall(pid);
			if err != nil {
				errChan <- err;
				log.Fatal(err);
			}	
			fmt.Printf("%v\n",res);
			p.Mu.Lock();
			info := p.Logs[tn];
			info.SYSCALL = res;
			p.Logs[tn] = info;
			p.Mu.Unlock();
			wg.Done();
		}(t);

		wg.Add(1);
		go func(tn time.Time) {
			res , err = ptracker.GetMem(pid);
			if err != nil {
				errChan <- err;
				log.Fatal(err);
			}	
			fmt.Printf("%v\n",res);
			p.Mu.Lock();
			info := p.Logs[tn];
			info.MEM = res;
			p.Logs[tn] = info;
			p.Mu.Unlock();
			wg.Done();
		}(t);
		
		wg.Wait();

		p.Mu.Lock()
        fmt.Printf("=== Completed trace %d at %v ===\n", len(p.Logs), t)
        p.Mu.Unlock()

        if(len(errChan) > 0){
			break;
		}

		time.Sleep(time.Second);
	}


	return fmt.Errorf("error occured in channel");
}