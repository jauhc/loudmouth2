package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/jauhc/go-csgsi"
	"github.com/ziutek/telnet"
)

// LoudSettings represents root level config: user, configs
type LoudSettings struct {
	User   string     `json:"user"` // your steamcommunity id
	Config LoudConfig `json:"config"`
}

// LoudConfig contains settings hardcoded until UI support?
type LoudConfig struct {
	State    bool `json:"state"`
	Clanid   bool `json:"clanid"`
	Clanfx   bool `json:"clanfx"`
	Owo      bool `json:"owo"`
	Kills    bool `json:"kills"`
	Deaths   bool `json:"deaths"`
	Greets   bool `json:"greets"`
	Ammowarn bool `json:"ammowarn"`
}

// read and apply user, configs etc
func readConfig(file string) LoudSettings {
	var settings LoudSettings
	settingsFile, err := os.Open(file)
	ec(err)
	defer settingsFile.Close()
	byteVal, err := ioutil.ReadAll(settingsFile)
	ec(err)
	json.Unmarshal(byteVal, &settings)
	return settings
}

// dum sleep wrapper because im lazy
func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// telnet writer
func consoleSend(t *telnet.Conn, s string) {
	t.Write([]byte(s + "\n"))
}

// check if given gun IS a gun
func isGun(gun *csgsi.Weapon) bool {
	// taser returns as ""
	if gun != nil {

		w := strings.ToLower(gun.Type)
		switch w {
		case "pistol":
			return true
		case "rifle":
			return true
		case "sniperrifle":
			return true
		case "submachine gun":
			return true
		case "machine gun":
			return true
		case "shotgun":
			return true
		default:
			return false
		}
	}
	return false
}

// forked a repo just to export structs lole
// gets given player's active weapon
func getActiveGun(gsi *csgsi.State) *csgsi.Weapon {
	for w := range gsi.Player.Weapons {
		if gsi.Player.Weapons[w].State == "active" {
			return gsi.Player.Weapons[w]
		}
	}
	return nil
}
