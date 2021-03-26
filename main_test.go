package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"nescio-server/pkg"
	"net"
	"testing"
	"time"
)

func TestTcpClient(t *testing.T) {
	c, err := net.Dial("tcp4", "127.0.0.1:3000")
	if err != nil {
		t.Fatal(err)
	}
	dec := json.NewDecoder(c)
	enc := json.NewEncoder(c)
	data := pkg.TcpData{Username: "fahim", Password: "fahim"}
	if err = enc.Encode(data); err != nil {
		t.Fatal(err)
	}
	if err = dec.Decode(&data); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", data)

	var b bytes.Buffer
	udp := pkg.UdpData{
		ID:        1000,
		Timestamp: 2000,
		Command:   4,
	}
	go func() {
		if err = binary.Write(&b, binary.BigEndian, &udp); err != nil {
			t.Fatal(err)
		}
		u, err := net.Dial("udp4", "127.0.0.1:3001")
		if err != nil {
			t.Fatal(err)
		}
		_, err = u.Write(b.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		buffer := make([]byte, 1024)
		n, _ := u.Read(buffer)
		data := bytes.NewBuffer(buffer[0:n])
		udpData := pkg.UdpData{}
		if err = binary.Read(data, binary.BigEndian, &udpData); err != nil {
			t.Log(err)
		}
		t.Logf("read from server with length %d: %v", n, udpData)
		<-time.After(1 * time.Second)
	}()
	for i := 0; i < 5; i++ {
		data = pkg.TcpData{Message: "ping"}
		if err = enc.Encode(data); err != nil {
			t.Fatal(err)
		}
		if err = dec.Decode(&data); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%v\n", data)
		<-time.After(3 * time.Second)
	}
	//<-time.After(10 * time.Second)
}
