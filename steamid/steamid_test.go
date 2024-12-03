package steamid_test

import (
	"github.com/bang-go/steam/steamid"
	"log"
	"testing"
)

func TestSteamId(t *testing.T) {
	//str := "STEAM_0:0:610610989"
	str := "76561199181487706"
	//str := "1221221978"
	//str := "[U:1:1221221978]"

	sid, err := steamid.New(str)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(sid.RenderSteamID2(false))
	log.Println(sid.RenderSteamID64())
	log.Println(sid.RenderSteamID3())
	log.Println(sid.GetAccountID())
	log.Println(sid.GetAccountType())
	log.Println(sid.GetInstance())
	log.Println(sid.GetUniverse())
	log.Println(sid.String())
}
