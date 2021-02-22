package main

import (
	"fmt"
	"os"
	"os/signal"
	_ "ottopoint-purchase/logging"
	"ottopoint-purchase/routers"
	"ottopoint-purchase/utils"
	"runtime"
	"strconv"
	"syscall"
)

func main() {
	maxpc, _ := strconv.Atoi(utils.GetEnv("MAXPROCS", "1"))
	runtime.GOMAXPROCS(maxpc)
	var errChan = make(chan error, 1)
	go func() {
		errChan <- routers.Server(utils.GetEnv("OTTOPOINT_PURCHASE", "0.0.0.0:8006"))
	}()
	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		fmt.Println("got an interrupt, exiting...")
	case err := <-errChan:
		if err != nil {
			fmt.Println("error while running api, exiting...", err)
		}
	}

}
