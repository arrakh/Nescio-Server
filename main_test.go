package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"nescio-server/pkg"
	"net"
	"testing"
)

func TestTcpClient(t *testing.T) {
	c, err := net.Dial("tcp4", "127.0.0.1:3000")
	if err != nil {
		t.Fatal(err)
	}

	/*go func() {
		var b bytes.Buffer
		udp := pkg.ClientData{
			ID:          1000,
			Timestamp:   2000,
			CommandEnum: 4,
		}

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
		udpData := pkg.ClientData{}
		if err = binary.Read(data, binary.BigEndian, &udpData); err != nil {
			t.Log(err)
		}
		t.Logf("read from server with length %d: %v", n, udpData)
		<-time.After(1 * time.Second)
	}()*/

	login, err := json.Marshal(pkg.Login{Username: "fahim", Password: "fahim"})
	buffer := new(bytes.Buffer)
	_ = binary.Write(buffer, binary.LittleEndian, pkg.ClientData{CommandEnum: 1})
	_ = binary.Write(buffer, binary.LittleEndian, login)

	if _, err := c.Write(buffer.Bytes()); err != nil {
		t.Fatal(err)
	}

	b := make([]byte, pkg.BUFFER)
	n, err := c.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("server data binary %v", b[0:n])
	var serverData pkg.ServerData
	_ = binary.Read(bytes.NewBuffer(b[0:n]), binary.LittleEndian, &serverData)
	t.Logf("server data %v", b[0:n])
	t.Logf("server data %v", serverData)

	/*for i := 0; i < 10; i++ {
		commandEnum := pkg.ClientData{CommandEnum: 0}
		buffer.Reset()
		_ = binary.Write(buffer, binary.LittleEndian, commandEnum)
		if _, err = c.Write(buffer.Bytes()); err != nil {
			t.Fatal(err)
		}
		<-time.After(3 * time.Second)
	}
	<-time.After(10 * time.Second)*/
}
