package pkg

import "net"

type Player struct {
	TcpConn *net.TCPConn
	UdpConn *net.UDPConn
	Profile Profile
}

type Profile struct {
	ID       uint32 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Success  bool   `json:"success,omitempty"`
}

type Vector2 struct {
	CoordinateX float64
	CoordinateY float64
	Stamp       int64
}
