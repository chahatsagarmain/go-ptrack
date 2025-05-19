package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chahatsagarmain/go-ptrack/internal/controller"
	"github.com/chahatsagarmain/go-ptrack/internal/process"
	"github.com/chahatsagarmain/go-ptrack/internal/writer"
	"github.com/urfave/cli/v3"
)

func main() {
    var customPath string
	var version string 
	version = "1.0.0" 
    cmd := &cli.Command{
        Name:  "ptrack",
        Usage: "Your low level process tracker in go",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:        "path",
                Usage:       "Custom path string (optional)",
                Destination: &customPath,
            },
        },
        Commands: []*cli.Command{
            {
                Name:  "track",
                Usage: "Track a process by PID",
                Action: func(ctx context.Context, cmd *cli.Command) error {
                    if cmd.NArg() == 1 {
                        pidRaw := cmd.Args().Get(0)
                        pid, err := strconv.ParseInt(pidRaw, 10, 64)
                        fmt.Printf("pid of the process is %v\n",pid);
                        if err != nil {
                            fmt.Printf("error parsing the arguement\n")
                            return fmt.Errorf("error parsing the arguement\n")
                        }
                        procPath := fmt.Sprintf("/proc/%d", pid)
                        fmt.Printf("path to write traces is : %v\n",customPath);
                        if _, err := os.Stat(procPath); os.IsNotExist(err) {
                            fmt.Printf("process with pid %d does not exist\n", pid)
                            return fmt.Errorf("process with pid %d does not exist\n", pid)
                        }
                        if customPath == "" {
							customPath = "/tmp/ptrack"
						}
                        fmt.Printf("tracing process %v now !\n", pid)
                        proc := process.NewProcess(int(pid))
                        done := make(chan error)

                        go func() {
                            for {
                                select {
                                case <-done:
                                    return
                                default:
                                    fmt.Println("tracing.....")
                                    proc.Mu.Lock()
                                    fmt.Printf("traces generated : %d\n", len(proc.Logs))
                                    proc.Mu.Unlock()
                                    if err := writer.WriteTrace(customPath , proc) ; err != nil {
                                        done <- err;
                                        return;
                                    }
                                    time.Sleep(5 * time.Second)
                                }
                            }
                        }()
                        err = controller.ControllerStart(proc.PID, 1000, proc)
                        close(done)
                        if err != nil {
                            log.Printf("tracing stopped....\n")
                            if err2 := writer.WriteTrace(customPath , proc) ; err2 != nil {
                                return err2;
                            }
                        }
                    } else {
                        fmt.Printf("enter the PID of prcess to trace\n")
                    }
                    return nil
                },
            },
            {
                Name:        "version",
                Usage:       "show version of ptracker",
				Action: func(ctx context.Context, c *cli.Command) error {
					fmt.Printf("version of ptracker is %v\n",version);
					return nil
				},
            },
        },
    }

    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatalf("closing ptrack : %v",err);
    }
}