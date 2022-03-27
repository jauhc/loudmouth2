package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/jauhc/go-csgsi"
)

/*---------legends--------\\
||    ! = IMPORTANT		  ||
||	? = possible feature  ||
||  % = to be improved    ||
\\------------------------//
	TODO
		! launch params
		!% logics
		? location parser with m_szLastPlaceName and "getpos"
		! error checks
	just rewrite shit tbh
*/

func killCheck(state *csgsi.State) {
	featureKillAnnounce(state)
}

// dont read this either thanks
func deathCheck(state *csgsi.State) {
	featureDeathAnnounce(state)
}

func gayTicker() {
	tick := 1
	// calls functions in features.go
	for range gayTimer.C {
		if tick >= 11 {
			tick = 1
		}

		if tick%8 == 0 { // every 800 ms
			featureRadioSpam()
		}

		if tick%9 == 0 { // every 900 ms
			featureSendChat()
		}

		if tick%5 == 0 { // every 500 ms
			featureClan()
		}
		tick++
	}
}

// main logic here with game events
func stateParser(gsi *csgsi.Game) {
	go func() {
		log.Println("starting parse..")
		createTimers()
		createConsoleCommands()
		for state := range gsi.Channel {
			if settings.Config.Ammowarn {
				featureAmmoWarning(&state) // warms when ammo is low :)
			}
			if stateOK(&state) { // sort out the "Previously" shit
				if state.Round.Phase == "live" && isLocalPlayer(&state) {
					if settings.Config.Kills {
						killCheck(&state)
					}
					if settings.Config.Deaths {
						deathCheck(&state)
					}
				}
			}
		}
	}()
}

func init() {
	seed := time.Now().UTC().UnixNano()
	seed ^= (seed << 12)
	seed ^= (seed >> 25)
	seed ^= (seed << 27)
	rand.Seed(seed) // big seed
	log.SetOutput(os.Stdout)

	if runtime.GOOS != "windows" {
		// can be removed if and when valve fixes crashing on linux
		println("invalid OS")
		os.Exit(1)
	}
}

func main() { // why the fuck nothing prints
	println("---START---")
	log.Println("Finding telnet creds...")
	getTelnetParams()
	log.Println("Generating 'hash'...")
	for {
		terribleHash = generateTerribleHash(9)
		if len(terribleHash) > 1 {
			break
		}
	}

	settings.User = getSteamID(true)
	// go startPanelServer() // web page access panel i never finished
	run("PASS " + settings.Pass)
	go listenerLoop(t) // thread for listening to rcon
	log.Println("Console connected!")
	gsi := createStateListener()
	sleep(130)
	log.Println("Listener created!")
	stateParser(gsi)

	gsi.Listen(fmt.Sprintf(":%s", settings.Gsiport)) // telnet port
	log.Println("----END----")
}
