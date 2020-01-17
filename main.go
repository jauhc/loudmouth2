package main

import (
	"log"
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

// fuck, redo; previous ONLY updates when worth updating
func killCheck(state *csgsi.State) {
	if state.Previously != nil {
		if state.Previously.Player != nil {
			if state.Previously.Player.Match_stats != nil {
				// i feel like im missing something very important
				if state.Previously.Player.Match_stats.Kills < state.Player.Match_stats.Kills {
					log.Println("KILL")
					log.Printf("curkills: %d", state.Player.Match_stats.Kills)
					log.Printf("prevkills: %d", state.Previously.Player.Match_stats.Kills)
					return
				}
			}
		}
	}
}

// check for added instead
func deathCheck(state *csgsi.State) {
	if state.Previously != nil {
		if state.Previously.Player != nil {
			if state.Previously.Player.Match_stats != nil {
				// i feel like im missing something very important
				if state.Previously.Player.Match_stats.Deaths < state.Player.Match_stats.Deaths &&
					state.Previously.Player.Match_stats.Deaths > 0 {
					log.Println("DETH")
					log.Printf("curdeth: %d", state.Player.Match_stats.Deaths)
					log.Printf("prevdeth: %d", state.Previously.Player.Match_stats.Deaths)
					return
				}
			}
		}
	}
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

var (
	beep      = syscall.MustLoadDLL("Kernel32.dll").MustFindProc("Beep")
	stateWait sync.Once
)

func main() {
	log.Println("---START---")
	t := creatListener()
	go listenerLoop(t) // thread for listening to rcon
	log.Println("Console connected!")
	gsi := createStateListener()
	log.Println("Listener created!")
	stateParser(gsi)
	gsi.Listen(":1489")
	log.Println("----END----")
}
