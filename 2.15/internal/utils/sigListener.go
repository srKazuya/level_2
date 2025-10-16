package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func SigListener() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		for sig := range sigs {
			if sig == syscall.SIGINT {
				fmt.Println()	
			}	
		}
	}()
}
