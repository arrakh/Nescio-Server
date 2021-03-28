package pkg

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"time"
)

func (player Player) HandleClient() {
	log.Debugf("serving %s\n", player.TcpConn.RemoteAddr().String())
	defer func() {
		log.WithField("id", player.ID).Debugln("tcp handler exited")
	}()

	authCh := make(chan Login)
	pingCh := make(chan bool)

	go func() {
		for {
			buffer := make([]byte, BUFFER)
			n, err := player.TcpConn.Read(buffer)
			if err != nil {
				player.errCh <- err
				continue
			}
			commandEnum := binary.LittleEndian.Uint32(buffer[0:4])
			log.WithField("player", player.Username).Debugf("command from player: %d", commandEnum)
			switch commandEnum {
			case PING:
				pingCh <- true
			case LOGIN:
				var login Login
				if err := json.Unmarshal(buffer[16:n], &login); err != nil {
					b, _ := json.Marshal(ClientData{CommandEnum: LOGIN})
					if _, err = player.TcpConn.Write(b); err != nil {
						player.errCh <- err
					}
				} else {
					authCh <- login
				}
			}
		}
	}()

	for {
		select {
		case err := <-player.errCh:
			log.WithField("player", player.Username).Debug(err)
			if netErr, ok := err.(net.Error); (ok && netErr.Timeout()) || err == io.EOF {
				_ = player.TcpConn.Close()
				connectedPlayer.Delete(player.ID)
				log.WithField("playerID", player.ID).Error("disconnected")
				return
			}
		case login := <-authCh:
			// TODO query from database or redis
			if (login.Username == "fahim" && login.Password == "fahim") ||
				(login.Username == "arrakh" && login.Password == "arrakh") {
				// give new id
				player.ID = uuid.New().ID()
				player.Username = login.Username
				log.WithField("player", player.Username).Println("login success")
				// set buffer
				buffer := new(bytes.Buffer)
				// write to buffer
				err := binary.Write(buffer, binary.LittleEndian, ServerData{CommandEnum: LOGIN})
				if err != nil {
					log.Debugln(err)
				}
				err = binary.Write(buffer, binary.LittleEndian, player.ID)
				if err != nil {
					log.Debugln(err)
				}
				if _, err = player.TcpConn.Write(buffer.Bytes()); err != nil {
					player.errCh <- err
				} else {
					var existingPlayerID uint32
					connectedPlayer.Range(func(key, value interface{}) bool {
						if value.(Player).Username == player.Username {
							log.WithField("player", player.Username).Debugf("found previous id %d\n", value.(Player).ID)
							existingPlayerID = value.(Player).ID
							return false
						}
						return true
					})
					if existingPlayerID != 0 {
						if existingPlayer, ok := connectedPlayer.Load(existingPlayerID); ok {
							log.WithField("player", player.Username).Debugf("disconnect previous id %d\n", existingPlayerID)
							existingPlayer.(Player).errCh <- io.EOF
						}
					}
					connectedPlayer.Store(player.ID, player)
				}
				player.setDeadline()
			} else {
				// set buffer
				buffer := new(bytes.Buffer)
				// write to buffer
				err := binary.Write(buffer, binary.LittleEndian, ServerData{CommandEnum: LOGIN})
				if err != nil {
					log.Debugln(err)
				}
				err = binary.Write(buffer, binary.LittleEndian, player.ID)
				if err != nil {
					log.Debugln(err)
				}
				if _, err = player.TcpConn.Write(buffer.Bytes()); err != nil {
					player.errCh <- err
				}
			}
		case <-pingCh:
			log.WithField("player", player.Username).Debugln("ping")
			buffer := new(bytes.Buffer)
			if err := binary.Write(buffer, binary.LittleEndian, ServerData{CommandEnum: PING}); err != nil {
				player.errCh <- err
			} else {
				log.WithField("player", player.Username).Debugln("pong")
				player.setDeadline()
			}
		}
	}
}

func (player Player) setDeadline() {
	if err := player.TcpConn.SetDeadline(time.Now().Add(DEADLINE * time.Second)); err != nil {
		player.errCh <- err
	}
}
