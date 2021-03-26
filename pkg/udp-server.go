package pkg

import (
	"fmt"
	shutdown "github.com/klauspost/shutdown2"
	log "github.com/sirupsen/logrus"
	"net"
)

var udpPort = ":3001"

func UDPServer() {
	udpAddr, err := net.ResolveUDPAddr("udp4", udpPort)
	if err != nil {
		log.Println("UDP server wrong address")
		shutdown.Exit(1)
	}

	udpConn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("UDP server starts %s\n", udpConn.LocalAddr().String())

	shut := shutdown.First()
	close(<-shut)
	log.Print("UDP server closed")
}
