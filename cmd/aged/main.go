package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/majestrate/aged"
)

func main() {

	log.Println("Starting up...")
	daemon, err := aged.NewDaemon()
	if err != nil {
		log.Fatal(err)
	}
	daemon.AgeProvider =
		&aged.LocalAgeProvider{
			DoBProvider: &aged.FlatFileDoBProvider{
				Filename: fmt.Sprintf("%s/.config/birthday", os.Getenv("HOME")),
			},
		}
	daemon.SetupHandler()
	aged.NotifyReady()
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGHUP)
	for {
		select {
		case sig := <-signalChan:
			if sig == os.Interrupt {
				aged.NotifyShutdown()
			}
			daemon.HandleSignal(sig)
			if !daemon.Running() {
				return
			}
		}
	}
}
