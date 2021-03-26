package pkg

import (
	shutdown "github.com/klauspost/shutdown2"
	log "github.com/sirupsen/logrus"
	"net"
)

var tcpPort = ":3000"

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

	tcpConnections := make(chan net.TCPConn)

	go func() {
		log.Printf("TCP server starts %s\n", tcpListener.Addr().String())
		for {
			c, err := tcpListener.AcceptTCP()
			if err != nil {
				log.Println(err)
			}
			tcpConnections <- *c
		}
	}()

	notifier := shutdown.Second()

	for {
		select {
		case conn := <-tcpConnections:
			player := Player{TcpConn: &conn}
			go player.HandleClient()
		case shut := <-notifier:
			close(shut)
			log.Print("TCP server closed")
		}
	}
}
