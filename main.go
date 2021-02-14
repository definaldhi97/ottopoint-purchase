package main

import (
	"fmt"
	"os"
	"os/signal"
	"ottopoint-purchase/routers"
	"syscall"
)

func main() {
	var errChan = make(chan error, 1)

	fmt.Println("test")

	// start router with designated port
	go func() {
		fmt.Println("Starting")
		errChan <- routers.Server(os.Getenv("OTTOPOINT_PURCHASE"))
	}()

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		fmt.Println("Got an interrupt, exiting...")
	case err := <-errChan:
		if err != nil {
			fmt.Println("Error while running api, exiting:", err)
		}
	}
}
