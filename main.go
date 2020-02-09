package main

import (
	"fmt"
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

// DONT READ THIS
func killCheck(state *csgsi.State) {
	if stateOK(state) {
		if state.Previously.Player.Match_stats.Kills < state.Player.Match_stats.Kills {
			//speechBuffer.PushBack(tellKill(state))
			run("enemydown")
			return
		}
	}
}

// dont read this either thanks
func deathCheck(state *csgsi.State) {
	if stateOK(state) {
		if state.Previously.Player.Match_stats.Deaths < state.Player.Match_stats.Deaths &&
			state.Previously.Player.Match_stats.Deaths > 0 {
			log.Println("DETH")
			// speechBuffer.PushBack(tellDeath(state))
			// TODO add to list instead
			return
		}
	}
}

// could be cleaner
func clanTicker() {
	output := fmt.Sprintf("cl_clanid 0")
	clanList := []int{7670261, 7670266, 7670268, 7670273, 7670276, 7670621, 7670634, 7670641, 7670647}
	for range clanTimer.C {
		if !settings.Config.Clanid && !settings.Config.Clanfx {
			if clanIdx > 0 {
				// run this once to clear clanid after disabling
				clanIdx = -1
				fxState = true
				run("cl_clanid 0")
			}
			return
		}
		if !settings.Config.Clanfx {
			if clanIdx >= len(clanList) {
				clanIdx = 0
			}
			clanIdx++
		} else if settings.Config.Clanfx {
			if clanIdx == 0 {
				fxState = true
			} else if clanIdx+1 >= len(clanList) {
				fxState = false
			}
			if fxState {
				clanIdx++
				// cl_clanid clanList[clanIdx++]
			} else if !fxState {
				clanIdx--
				// cl_clanid clanList[clanIdx--]
			}
		}
		output = fmt.Sprintf("cl_clanid %d", clanList[clanIdx])
		run(output)
	}
}

// radio spammer
func radioTicker() {
	for range radioTimer.C {
		if settings.Config.Radiospam {
			run("ohn")
		}
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
