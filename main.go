package main

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"os"
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
		% UI
		!% logics
		? location parser with m_szLastPlaceName and "getpos"
		! error checks
	just rewrite shit tbh
*/

// short for error check
func ec(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
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

// loop print messages
func listenerLoop(t *telnet.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := t.Read(buf)
		os.Stdout.Write(buf[:n])
		ec(err)
	}
}

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
	// why does it check currently spectating's stats
	if state.Previously != nil {
		if state.Previously.Player != nil {
			if state.Previously.Player.Match_stats != nil {
				// i feel like im missing something very important
				if state.Previously.Player.Match_stats.Kills < state.Player.Match_stats.Kills {
					//speechBuffer.PushBack(tellKill(state))
					run("enemydown")
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
					speechBuffer.PushBack(tellDeath(state))
					// TODO add to list instead
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

	return "bazinga\nenemydown"
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
	// return cheese[rand.Intn(len(cheese))]
	_ = cheese
	return "oh no"
}

// could be cleaner
func clanTicker() {
	output := fmt.Sprintf("cl_clanid 0")
	clanList := []int{7670261, 7670266, 7670268, 7670273, 7670276, 7670621, 7670634, 7670641, 7670647}
	for range clanTimer.C {
		if !settings.Config.Clanid {
			return
		}
		if !settings.Config.Clanfx {
			if clanIdx >= len(clanList) {
				clanIdx = 0
			}
			output = fmt.Sprintf("cl_clanid %d", clanList[clanIdx])
			clanIdx++
		} else if settings.Config.Clanfx {
			if clanIdx == 0 {
				fxState = true
			} else if clanIdx+1 >= len(clanList) {
				fxState = false
			}
			if fxState {
				output = fmt.Sprintf("cl_clanid %d", clanList[clanIdx])
				clanIdx++
				// cl_clanid clanList[clanIdx++]
			} else if !fxState {
				output = fmt.Sprintf("cl_clanid %d", clanList[clanIdx])
				clanIdx--
				// cl_clanid clanList[clanIdx--]
			}
		}
		run(output)
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

// main logic here with game events
func stateParser(gsi *csgsi.Game) {
	go func() {
		log.Println("starting parse..")
		// go speechTicker()
		go clanTicker()
		for state := range gsi.Channel {
			if state.Round.Phase == "live" || 1 == 1 {

				// local player check
				if state.Player.SteamId == settings.User {
					// log.Println("live")
					//TODO settings checks
					ammoWarning(&state) // warms when ammo is low :)
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

func run(cheese string) {
	output := fmt.Sprintf("%s\n", cheese)
	t.Write([]byte(output))
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
	sleep(500)
	log.Println("Listener created!")
	stateParser(gsi)

	gsi.Listen(":1489")
	log.Println("----END----")
}
