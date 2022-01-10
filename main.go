package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ctx, cancel := context.WithCancel(context.Background())

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sigterm:
		log.Println("receive stop signal")
	}

	cancel()
	time.Sleep(3 * time.Second)
}
