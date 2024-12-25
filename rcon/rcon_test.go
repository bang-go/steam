package rcon_test

import (
	"github.com/bang-go/steam/rcon"
	"log"
	"testing"
	"time"
)

func TestRcon(t *testing.T) {
	var err error
	rc := rcon.New(&rcon.Config{Addr: "127.0.0.1:27015", Timeout: time.Second * 5})
	err = rc.Dail()
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()
	err = rc.Auth("123456")
	if err != nil {
		log.Fatal(err)
	}
	data, err := rc.ExecCommand([]byte("bot_add_ct"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(data)
}
