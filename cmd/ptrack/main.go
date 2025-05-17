package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name: "ptrack",
		Usage: "Your low level process tracker in go",
		Action: func(ctx context.Context , cmd *cli.Command) error {
			if(cmd.NArg() == 1){
				pidRaw := cmd.Args().Get(0);
				pid , err := strconv.ParseInt(pidRaw , 10 , 64);
				if err != nil {
					log.Fatalf("ERROR PARSING THE ARGUEMENT\n");
					return fmt.Errorf("ERROR PARSING THE ARGUEMENT\n");
				}
				procPath := fmt.Sprintf("/proc/%d", pid)
    			if _, err := os.Stat(procPath); os.IsNotExist(err) {
    			    log.Fatalf("PROCESS WITH PID %d DOES NOT EXIST\n", pid)
    			    return fmt.Errorf("PROCESS WITH PID %d DOES NOT EXIST\n", pid)
    			}
				fmt.Printf("tracing process %v now !",pid);
			} else if (cmd.NArg() > 1) {
				fmt.Printf("ENTER ONLY PID OF PROCESS TO TRACE\n");
			} else {
				fmt.Printf("ENTER THE PID OF PROCESS TO TRACE\n");
			}
			return nil;
		},
	}

	if err := cmd.Run(context.Background(),os.Args) ; err != nil {
		log.Fatalf("Error starting up ptrack")
	}
}