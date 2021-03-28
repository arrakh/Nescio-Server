package pkg

import (
	shutdown "github.com/klauspost/shutdown2"
	log "github.com/sirupsen/logrus"
	"net"
)

func TCPServer() {
	tcpAddress, err := net.ResolveTCPAddr("tcp4", TCP_ADDRESS)
	if err != nil {
		log.Println("TCP server wrong address format")
		shutdown.Exit(1)
	}

	tcpListener, err := net.ListenTCP("tcp4", tcpAddress)
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
