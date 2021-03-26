package pkg

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
	"time"
)

func (player Player) HandleClient() {
	log.Debugf("serving %s\n", player.TcpConn.RemoteAddr().String())
	defer func() {
		log.WithField("id", player.ID).Debugln("tcp handler exited")
	}()

	player.setDeadline()

	decoder := json.NewDecoder(player.TcpConn)
	encoder := json.NewEncoder(player.TcpConn)

	authCh := make(chan TcpData)
	pingCh := make(chan bool)

	go func() {
		for {
			var tcpData TcpData

			err := decoder.Decode(&tcpData)
			if err != nil {
				player.errCh <- err
			} else if len(tcpData.Username) > 0 && len(tcpData.Password) > 0 {
				authCh <- tcpData
			} else if len(tcpData.Message) > 0 && strings.ToLower(tcpData.Message) == "ping" {
				pingCh <- true
			}
		}
	}()

	for {
		select {
		case err := <-player.errCh:
			log.WithField("player", player.Username).Debug(err)
			if netErr, ok := err.(net.Error); (ok && netErr.Timeout()) || err == io.EOF {
				_ = player.TcpConn.Close()
				connectedPlayer.Delete(player.Username)
				log.WithField("player", player.Username).Error("disconnected")
				return
			}
		case authenticated := <-authCh:
			// TODO query from database or redis
			if authenticated.Username == "fahim" && authenticated.Password == "fahim" {
				player.ID = uuid.New().ID()
				player.Username = authenticated.Username
				log.WithField("player", player.Username).Println("login success")
				if err := encoder.Encode(TcpData{Username: player.Username, ID: player.ID}); err != nil {
					player.errCh <- io.EOF
					continue
				}
				if existingPlayer, ok := connectedPlayer.LoadOrStore(player.Username, player); ok {
					log.WithField("player", player.Username).Printf("disconnect previous id %d\n", existingPlayer.(Player).ID)
					_ = existingPlayer.(Player).TcpConn.Close()
					log.WithField("player", player.Username).Debugf("%v\n", existingPlayer.(Player).UdpConn)
					//_ = existingPlayer.(Player).UdpConn.Close()
					existingPlayer.(Player).errCh <- io.EOF
					connectedPlayer.Store(player.Username, player)
				}
				player.setDeadline()
			} else {
				player.errCh <- io.EOF
			}
		case <-pingCh:
			log.WithField("player", player.Username).Debugln("ping")
			if err := encoder.Encode(TcpData{Message: "pong"}); err != nil {
				player.errCh <- err
			} else {
				log.WithField("player", player.Username).Debugln("pong")
				player.setDeadline()
			}
		}
	}
}

func (player Player) setDeadline() {
	if err := player.TcpConn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		player.errCh <- err
	}
}
