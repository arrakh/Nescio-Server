package pkg

import (
	shutdown "github.com/klauspost/shutdown2"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

var tcpPort = ":3000"

var connectedPlayer sync.Map

func TCPServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpPort)
	if err != nil {
		log.Println("TCP server wrong address")
		shutdown.Exit(1)
	}

	tcpListener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Println(err)
	}

	tcpConnCh := make(chan net.TCPConn)

	go func() {
		log.Printf("TCP server starts %s\n", tcpListener.Addr().String())
		for {
			c, err := tcpListener.AcceptTCP()
			if err != nil {
				log.Println(err)
			}
			tcpConnCh <- *c
		}
	}()

	notifier := shutdown.Second()

	for {
		select {
		case tcpConn := <-tcpConnCh:
			go NewPlayer(tcpConn).HandleClient()
		case shut := <-notifier:
			_ = tcpListener.Close()
			close(tcpConnCh)
			close(shut)
			log.Print("TCP server closed")
			return
		}
	}
}
