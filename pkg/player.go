package pkg

import (
	"net"
	"sync"
)

var connectedPlayer sync.Map

// TODO will be put inside config.ini file
const (
	BUFFER      = 1024
	TCP_ADDRESS = ":3000"
	UDP_ADDRESS = ":3001"
	DEADLINE    = 30
)

type Player struct {
	TcpConn   *net.TCPConn
	UdpConn   *net.UDPConn
	ID        uint32
	Username  string
	errCh     chan error
	udpDataCh chan []byte
}

type Login struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

const (
	PING = iota
	LOGIN
	CHAR_MOVE
	CHAR_ATTACK
)

type ClientData struct {
	CommandEnum uint32
	ID          uint32
	Timestamp   int64
}

type ServerData struct {
	CommandEnum uint32
	Counter     int32
}

type PingPongData struct {
	Ping int32
}

type Vector2 struct {
	CoordinateX float64
	CoordinateY float64
	Stamp       int64
}

func NewPlayer(tcpConn net.TCPConn) Player {
	return Player{
		TcpConn: &tcpConn,
		errCh:   make(chan error),
	}
}
