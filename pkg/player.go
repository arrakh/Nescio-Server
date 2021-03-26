package pkg

import (
	"net"
)

type Player struct {
	TcpConn   *net.TCPConn
	UdpConn   *net.UDPConn
	ID        uint32
	Username  string
	errCh     chan error
	udpDataCh chan []byte
}

type TcpData struct {
	ID       uint32 `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Success  bool   `json:"success,omitempty"`
	Message  string `json:"message,omitempty"`
}

type UdpData struct {
	ID        uint32
	Timestamp int64
	Command   int32
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
