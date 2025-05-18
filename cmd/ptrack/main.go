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
			&cli.StringFlag{
                Name:        "version",
                Usage:       "show version of ptracker",
				Action: func(ctx context.Context, c *cli.Command, s string) error {
					fmt.Printf("version of ptracker is %v",version);
					return nil
				},
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
                        if err != nil {
                            log.Fatalf("error parsing the arguement\n")
                            return fmt.Errorf("error parsing the arguement\n")
                        }
                        procPath := fmt.Sprintf("/proc/%d", pid)
                        if customPath != "" {
                            procPath = customPath
                        } else {
							customPath = "/tmp/ptrack"
						}
                        if _, err := os.Stat(procPath); os.IsNotExist(err) {
                            log.Fatalf("process with pid %d does not exist\n", pid)
                            return fmt.Errorf("process with pid %d does not exist\n", pid)
                        }
                        fmt.Printf("tracing process %v now !\n", pid)
                        proc := process.NewProcess(int(pid))
                        done := make(chan struct{})

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
                                    time.Sleep(5 * time.Second)
                                }
                            }
                        }()
                        err = controller.ControllerStart(proc.PID, 1000, proc)
                        close(done)
                        if err != nil {
                            log.Printf("tracing stopped....\n")
                            log.Printf("writing traces to path %v\n...",customPath);
							if _ , err := os.Stat(customPath) ; err != nil && customPath != "/tmp/ptrack" {
								log.Printf("path specified does not exist\n")
								return err
							}
							json , err2 := proc.ToJSON();
							if err2 != nil{
								log.Printf("error parsings logs\n");
								return err2
							}
							if customPath == "/tmp/ptrack" {
								err := os.MkdirAll(customPath, 0777)
								if err != nil {
								    log.Fatalf("failed to create /tmp/ptrack: %v", err)
								}
								if err := os.WriteFile(fmt.Sprintf("%v/ptrack.json",customPath),json,0666) ; err != nil {
									log.Printf("error writing file...\n");
									return err;
								}
								log.Println("successfully wrote to path");
							} else {
								if err := os.WriteFile(customPath,json,0666) ; err != nil {
									log.Printf("error writing file...\n");
									return err;
								}
								log.Println("successfully wrote to path");
							}
							return err
                        }
                    } else {
                        fmt.Printf("enter the PID of prcess to trace\n")
                    }
                    return nil
                },
            },
        },
    }

    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatalf("closing ptrack : %v",err);
    }
}