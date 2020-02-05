package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jauhc/go-csgsi"
	"github.com/ziutek/telnet"
)

// LoudSettings represents root level config: user, configs
type LoudSettings struct {
	User   string     `json:"user"` // your steamcommunity id
	Pass   string     `json:"pass"` // console connection password, same as -netconpassword
	Config LoudConfig `json:"config"`
}

// password is static until i can figure out how to automate this
// shit without forcing this tool to be a csgo launcher

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

const ()

var (
	// clan shit
	fxState = true // true fowards, false backwards
	clanIdx = 0

	// tickers (timers) to handle spamming
	speechTimer = time.NewTicker(900 * time.Millisecond)
	clanTimer   = time.NewTicker(500 * time.Millisecond)

	speechBuffer = list.New()

	// beeping for low ammo warning //TODO add OS check)
	beep      = syscall.MustLoadDLL("Kernel32.dll").MustFindProc("Beep")
	stateWait sync.Once

	t        = creatListener()
	settings = readConfig("loud.json")
)

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

// short for error check
func ec(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

// loop print messages
func listenerLoop(t *telnet.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := t.Read(buf)
		os.Stdout.Write(buf[:n])
		ec(err)
	}
}

func speechTicker() {
	// TODO
	// make an instant trigger then check if buffer > 0 and go here
	// or; make a 1 sec timer which toggles bool whenever buffer needed
	// and fire said timer when sent message with no buffer
	for range speechTimer.C {
		if speechBuffer.Len() > 0 {
			pop := speechBuffer.Front()
			poop := fmt.Sprintf("say %s", pop.Value)
			run(poop)
			speechBuffer.Remove(pop)
		}
	}
}

// creates telnet client
func creatListener() *telnet.Conn {
	t, err := telnet.Dial("tcp", ":2121")
	ec(err)
	log.Println("dialed")
	os.Stdin.WriteString("poop") // what
	return t
}

// creates gsi listener
func createStateListener() *csgsi.Game {
	game := csgsi.New(3)
	return game
}

// run console command
func run(cheese string) {
	t.Write([]byte(cheese + "\n"))
}

func isLocalPlayer(state *csgsi.State) bool {
	return state.Player.SteamId == settings.User
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
