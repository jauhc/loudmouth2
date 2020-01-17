package main

import (
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

/*---------legends--------\\
||    ! = IMPORTANT		  ||
||	? = possible feature  ||
||  % = to be improved    ||
\\------------------------//
	TODO
		! launch params
		! http server for listening to game
		% UI
		!% logics
		% wrapper for timers
		? location parser with m_szLastPlaceName and "getpos"
		! error checks
	just rewrite shit tbh

	options
		gsi https://github.com/dank/go-csgsi
		telnet https://github.com/ziutek/telnet/

	flow
		attempt to connect (telnet prio)
		check configs -> read/write them
		go to listen loop
*/

// short TODO
// how to wait for goroutines to exit before main quits
// redirect their output to main's

// short for error check
func ec(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

// dum sleep wrapper because im lazy
func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// creates telnet client
func creatListener() *telnet.Conn {
	t, err := telnet.Dial("tcp", ":2121")
	ec(err)
	log.Println("dialed")
	os.Stdin.WriteString("poop")
	return t
}

// creates gsi listener
func createStateListener() *csgsi.Game {
	game := csgsi.New(3)
	return game
}

// telnet writer
func consoleSend(t *telnet.Conn, s string) {
	t.Write([]byte(s + "\n"))
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

func ammoWarning(state *csgsi.State) {
	w := getActiveGun(state)
	if isGun(w) {
		if float32(w.Ammo_clip)/float32(w.Ammo_clip_max) < 0.3 &&
			w.Ammo_clip_max > 1 {
			// make check if we want this on or not
			beep.Call(80, 168)
		}
	}
}

// DONT READ THIS
func killCheck(state *csgsi.State) {
	if state.Previously != nil {
		if state.Previously.Player != nil {
			if state.Previously.Player.Match_stats != nil {
				// i feel like im missing something very important
				if state.Previously.Player.Match_stats.Kills < state.Player.Match_stats.Kills {
					say(tellKill(state))
					return
				}
			}
		}
	}
}

// dont read this either thanks
func deathCheck(state *csgsi.State) {
	if state.Previously != nil {
		if state.Previously.Player != nil {
			if state.Previously.Player.Match_stats != nil {
				// i feel like im missing something very important
				if state.Previously.Player.Match_stats.Deaths < state.Player.Match_stats.Deaths &&
					state.Previously.Player.Match_stats.Deaths > 0 {
					log.Println("DETH")
					return
				}
			}
		}
	}
}

func tellKill(state *csgsi.State) string {
	var cheese []string
	if state.Player.State.Flashed > 50 {
		cheese = []string{
			"hey i was blind",
			"owned while flashed",
			"ez blind kills",
			"i am blind, literally"}
	} else if state.Player.State.Smoked > 50 {
		cheese = []string{
			"~ from within the smoke ~",
			"puff puff im in the smokes",
			"really cloudy here",
			"i could barely see anything here wtf"}
	} else {
		cheese = []string{
			"blap blap",
			"sit down",
			"later",
			"hey, how about a break?",
			"you alright?",
			"hit or miss? guess i never miss, huh?",
			"ez",
			"ezpz",
			"you just got dabbed on!",
			"owned",
			"ownd",
			"whats happening with you",
			"get pooped on"}
	}

	postfix := []string{
		"kid",
		"kiddo",
		"nerd",
		"geek",
		"noob"}

	picked := cheese[rand.Intn(len(cheese))]

	if rand.Float32() > 0.6 {
		return fmt.Sprintf("%s %s\nenemydown", picked, postfix[rand.Intn(len(postfix))])
	}

	return fmt.Sprintf("%s\nenemydown", picked)
}

func tellDeath(state *csgsi.State) string {
	var cheese []string
	if state.Player.State.Flashed > 50 {
		cheese = []string{
			"i was blind lole",
			"how do i shoot blind?",
			"oops i was flashed",
			"help i cant see",
			"why is my screen white"}
	} else {
		cheese = []string{
			"oops",
			"i meant to do that :)",
			"wtf lag",
			"i was looking at the map",
			"excuse me?",
			"oh",
			"fricking tickrate",
			"omg 64 tick"}
	}
	return cheese[rand.Intn(len(cheese))]
}

// main logic here with game events
func stateParser(gsi *csgsi.Game) {
	go func() {
		for state := range gsi.Channel {
			if state.Round.Phase == "live" || 1 == 1 {
				// log.Println("live")
				// logics here
				ammoWarning(&state) // warms when ammo is low :)
				killCheck(&state)
				deathCheck(&state)
			}
		}
	}()
}

func say(cheese string) {
	output := fmt.Sprintf("say %s\n", cheese)
	t.Write([]byte(output))
}

var (
	beep      = syscall.MustLoadDLL("Kernel32.dll").MustFindProc("Beep")
	stateWait sync.Once
	t         = creatListener()
)

func init() {
	seed := time.Now().UTC().UnixNano()
	seed ^= (seed << 12)
	seed ^= (seed >> 25)
	seed ^= (seed << 27)
	rand.Seed(seed) // big seed
}

func main() {
	log.Println("---START---")
	go listenerLoop(t) // thread for listening to rcon
	log.Println("Console connected!")
	gsi := createStateListener()
	log.Println("Listener created!")
	stateParser(gsi)
	gsi.Listen(":1489")
	log.Println("----END----")
}
