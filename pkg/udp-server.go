package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	shutdown "github.com/klauspost/shutdown2"
	log "github.com/sirupsen/logrus"
	"net"
)

func UDPServer() {
	udpAddr, err := net.ResolveUDPAddr("udp4", UDP_ADDRESS)
	if err != nil {
		log.Println("UDP server wrong address format")
		shutdown.Exit(1)
	}

	udpListener, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		fmt.Println(err)
	}

	udpConnCh := make(chan *net.UDPAddr)
	udpDataCh := make(chan []byte)

	go func() {
		log.Printf("UDP server starts %s\n", udpListener.LocalAddr().String())
		buffer := make([]byte, 1024)
		for {
			n, udpClient, err := udpListener.ReadFromUDP(buffer)
			if err != nil {
				log.Debugln(err)
			}
			udpData := ClientData{}
			data := bytes.NewBuffer(buffer[0:n])
			log.Debugf("read from %s with length %d: %d", udpClient.String(), n, len(data.String()))
			if err = binary.Read(data, binary.BigEndian, &udpData); err != nil {
				log.Debugln(err)
				continue
			}
			log.Debugln(udpData)
			udpConnCh <- udpClient
			udpDataCh <- buffer[0:n]
		}
	}()

	notifier := shutdown.Second()

	for {
		select {
		case udpClient := <-udpConnCh:
			go func() {
				data := <-udpDataCh
				log.Debugf("write to addr %s: %s\n", udpClient.String(), string(data))
				udpListener.WriteToUDP(data, udpClient)
			}()
		case shut := <-notifier:
			_ = udpListener.Close()
			close(udpConnCh)
			close(shut)
			log.Print("UDP server closed")
			return
		}
	}
}
