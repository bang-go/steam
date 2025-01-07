package steamapi_test

import (
	"github.com/bang-go/steam/steamapi"
	"log"
	"testing"
)

func TestAccount(t *testing.T) {
	var err error
	apiKey := ""
	g := steamapi.NewGameServersService(&steamapi.GameServersServiceConfig{ApiKey: apiKey})
	resp, err := g.GetServerSteamIDsByIP([]string{"x.x.x.x:27015"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*resp)
}
