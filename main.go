package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fhofherr/rsched/internal/cmd"
	"github.com/fhofherr/rsched/internal/restic"
)

func main() {
	cfg, err := cmd.LoadConfig(os.Args[1:])
	if err != nil {
		fmt.Printf("\n%v\n", err)
		os.Exit(1)
	}
	rsched := &cmd.RSched{
		Scheduler: &restic.Scheduler{},
	}
	onSignal(rsched.Shutdown, syscall.SIGINT, syscall.SIGTERM)
	rsched.Run(cfg)
}

func onSignal(f func(), sigs ...os.Signal) {
	sigc := make(chan os.Signal, 1)

	go func() {
		<-sigc
		f()
	}()

	signal.Notify(sigc, sigs...)
}
