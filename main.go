package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/dank/go-csgsi"

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
func isGun(gsi *csgsi.State) bool {
	// taser returns as ""
	w := strings.ToLower(getActiveGunType(gsi))
	switch w {
	case "pistol":
		return true
	case "rifle":
		return true
	case "sniperrifle":
		return true
	case "machine gun":
		return true
	case "submachine gun":
		return true
	case "shotgun":
		return true
	default:
		return false
	}
}

// gets given player's active weapon
func getActiveGunType(gsi *csgsi.State) string {
	for w := range gsi.Player.Weapons {
		if gsi.Player.Weapons[w].State == "active" {
			return gsi.Player.Weapons[w].Type
		}
	}
	return ""
}

// would be better to just return weapon struct in one method?

// gets given player's active weapon
func getActiveGunName(gsi *csgsi.State) string {
	for w := range gsi.Player.Weapons {
		if gsi.Player.Weapons[w].State == "active" {
			return gsi.Player.Weapons[w].Name
		}
	}
	return ""
}

// main logic here with game events
func stateParser(gsi *csgsi.Game) {
	go func() {
		for state := range gsi.Channel {
			log.Println(state.Player.Name)
			if state.Round.Phase == "live" || 1 == 1 {
				log.Println("live")
				// logics here
				log.Println(getActiveGunType(&state))
			}
		}
	}()
}

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
