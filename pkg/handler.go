package pkg

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

func (player Player) HandleClient() {
	log.Printf("serving %s\n", player.TcpConn.RemoteAddr().String())

	for {
		d := json.NewDecoder(player.TcpConn)
		e := json.NewEncoder(player.TcpConn)

		var profile Profile

		err := d.Decode(&profile)
		if err != nil {
			log.Println(err)
			if err == io.EOF {
				log.Printf("closed %s\n", player.TcpConn.RemoteAddr().String())
				_ = player.TcpConn.Close()
			}
			return
		} else {
			<-time.After(5 * time.Second)
			if profile.Username == "fahim" && profile.Password == "fahim" {
				player.Profile.ID = uuid.New().ID()
				player.Profile.Username = profile.Username
				player.Profile.Success = true
				log.Println(player.Profile)
				e.Encode(player.Profile)
			}
		}
	}
}
