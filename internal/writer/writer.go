package writer

import (
	"fmt"
	"log"
	"os"

	"github.com/chahatsagarmain/go-ptrack/internal/process"
)

func writeFile(customPath string, json []byte) error {
	if customPath == "/tmp/ptrack" {
		err := os.MkdirAll(customPath, 0777)
		if err != nil {
			log.Fatalf("failed to create /tmp/ptrack: %v", err)
		}
		if err := os.WriteFile(fmt.Sprintf("%v/ptrack.json", customPath), json, 0666); err != nil {
			log.Printf("error writing file...\n")
			return err
		}
		log.Println("successfully wrote to path")
	} else {
		if err := os.WriteFile(customPath, json, 0666); err != nil {
			log.Printf("error writing file... : %v\n",err);
			return err
		}
		log.Println("successfully wrote to path")
	}

	return nil
}

func WriteTrace(customPath string, proc *process.Process) error {
	log.Printf("writing traces to path %v\n...", customPath)
	if _, err := os.Stat(customPath); err != nil && customPath != "/tmp/ptrack" {
		log.Printf("path specified does not exist\n")
		return err
	}
	json, err2 := proc.ToJSON()
	if err2 != nil {
		log.Printf("error parsings logs\n")
		return err2
	}
	if err2 := writeFile(customPath, json); err2 != nil {
		return err2
	}
	return nil
}
