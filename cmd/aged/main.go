package main

import (
	"fmt"
	"github.com/majestrate/aged"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	log.Println("Starting up...")
	daemon, err := aged.NewDaemon()
	if err != nil {
		log.Fatal(err)
	}
	daemon.AgeProvider = &aged.FlatAgeFile{
		Filename: fmt.Sprintf("%s/.config/birthday", os.Getenv("HOME")),
	}
	daemon.SetupHandler()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGHUP)
	select {
	case sig := <-signalChan:
		daemon.HandleSignal(sig)
		if !daemon.Running() {
			return
		}
	}
}
