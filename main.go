package main

import (
	"log"
	"math/rand"
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
		% UI
		!% logics
		? location parser with m_szLastPlaceName and "getpos"
		! error checks
	just rewrite shit tbh
*/

func ammoWarning(state *csgsi.State) {
	w := getActiveGun(state)
	if isGun(w) {
		if float32(w.Ammo_clip)/float32(w.Ammo_clip_max) < 0.3 &&
			w.Ammo_clip_max > 1 {
			if settings.Config.Ammowarn {
				beep.Call(80, 168)
			}
		}
	}
}

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

		if tick%9 == 0 {
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
			if stateOK(&state) {
				if state.Round.Phase == "live" && isLocalPlayer(&state) {
					if settings.Config.Ammowarn {
						ammoWarning(&state) // warms when ammo is low :)
					}
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
}

func main() {
	log.Println("---START---")
	run("PASS " + settings.Pass)
	go listenerLoop(t) // thread for listening to rcon
	log.Println("Console connected!")
	gsi := createStateListener()
	sleep(130)
	log.Println("Listener created!")
	stateParser(gsi)

	gsi.Listen(":1489")
	log.Println("----END----")
}
