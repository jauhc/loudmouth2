package main

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jauhc/go-csgsi"
	"github.com/ziutek/telnet"
)

const (
	configFile = "loud.json"
)

var (
	// clan shit
	fxState = true // true fowards, false backwards
	clanIdx = 0

	// tickers (timers) to handle spamming
	speechTimer = time.NewTicker(900 * time.Millisecond)
	clanTimer   = time.NewTicker(500 * time.Millisecond)
	radioTimer  = time.NewTicker(800 * time.Millisecond)

	speechBuffer = list.New()

	// beeping for low ammo warning //TODO add OS check)
	beep      = syscall.MustLoadDLL("Kernel32.dll").MustFindProc("Beep")
	stateWait sync.Once

	t            = creatListener()
	settings     = readConfig(configFile)
	terribleHash = generateTerribleHash(17)
)

func createTimers() {

	go radioTicker()
	// go speechTicker()
	go clanTicker()
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
		data := buf[:n]
		os.Stdout.Write(data)
		go consoleParse(data)
		ec(err)
	}
}

// not 100% sure if this works but
func say(cheese string) {
	speechBuffer.PushBack(cheese)
	if speechBuffer.Len() == 1 { // if just one, else let ticker handle it
		pop := speechBuffer.Front()
		output := fmt.Sprintf("say %s", pop.Value)
		run(output)
		speechBuffer.Remove(pop)
	}
}

// timer for chat output
func speechTicker() {
	for range speechTimer.C {
		if speechBuffer.Len() > 1 {
			pop := speechBuffer.Front()
			poop := fmt.Sprintf("say %s", pop.Value)
			run(poop)
			speechBuffer.Remove(pop)
		}
	}
}

// creates telnet client
func creatListener() *telnet.Conn {
start:
	t, err := telnet.Dial("tcp", ":2121")
	if err != nil {
		// no clue how to retry so we make this
		log.Println(err)
		goto start
	}
	log.Println("dialed")
	os.Stdin.WriteString("poop") // what
	return t
}

// creates gsi listener
func createStateListener() *csgsi.Game {
	game := csgsi.New(3)
	return game
}

// dumb wrapper for checking state structure integrity because i cant think of a better way
func stateOK(state *csgsi.State) bool {
	if state.Previously != nil {
		if state.Previously.Player != nil {
			if state.Previously.Player.Match_stats != nil {
				return true
			}
		}
	}
	return false
}

// run console command
func run(cheese string) {
	// fmt.Println("> " + cheese)
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

func consoleParse(data []byte) {
	toParse := string(data)
	if len(toParse) > 0 {
		hashIdx := strings.Index(toParse, terribleHash)
		if hashIdx > -1 {
			checkCvars(strings.Split(toParse[:hashIdx], " "))
		}
	}
}

func generateTerribleHash(howlong int) string {
	// charset := "iIl1|!o0OS5B8"
	charset := "iIl1|"
	b := make([]byte, howlong)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
