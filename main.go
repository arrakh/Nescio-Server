package main

import (
	shutdown "github.com/klauspost/shutdown2"
	log "github.com/sirupsen/logrus"
	"nescio-server/pkg"
	"os"
	"syscall"
	"time"
)

func main() {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	go pkg.UDPServer()
	go pkg.TCPServer()

	shutdown.OnSignal(0, os.Interrupt, syscall.SIGTERM)
	shutdown.SetTimeout(time.Second * 10)
	shut := shutdown.Second()
	<-shut
	log.Info("Exited")
}
