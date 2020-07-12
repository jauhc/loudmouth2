package main

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jauhc/go-csgsi"
)

const (
	configFile = "loud.json"
	msgCode    = "‎ : " // DONT TOUCH OR WE ALL DIE
	uniqueCode = "‎"    // DONT TOUCH OR WE ALL DIE
)

var (
	// clan shit
	fxState = true // true fowards, false backwards
	clanIdx = 0

	// not actually sure if best approach
	gayTimer = time.NewTicker(100 * time.Millisecond)

	speechBuffer = list.New()

	// beeping for low ammo warning
	// no need to worry about linux since it (the game) crashes on connection
	beep      = syscall.MustLoadDLL("Kernel32.dll").MustFindProc("Beep")
	stateWait sync.Once

	t            = creatListener(":2121")
	settings     = readConfig(configFile)
	terribleHash = generateTerribleHash(9)
)

func createTimers() {
	go gayTicker()
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
func listenerLoop(sock net.Conn) {
	for {
		recvData := make([]byte, 256)
		n, err := sock.Read(recvData)

		ec(err)
		if n > 0 {
			out := string(recvData)
			idx := strings.IndexByte(out, '\x00')
			if idx > -1 {
				os.Stdout.Write(recvData[:idx])
				consoleParse(out)
			} else {
				os.Stdout.Write(recvData)
				consoleParse(out)
			}
		}
	}
}

// not 100% sure if this works but
func say(cheese string, isTeam ...bool) {
	var output string
	speechBuffer.PushBack(cheese)
	if speechBuffer.Len() == 1 { // if nothing else in queue, dont queue it
		pop := speechBuffer.Front()
		if len(isTeam) > 0 {
			if isTeam[0] {
				output = fmt.Sprintf("say_team %s", pop.Value)
			}
		} else {
			output = fmt.Sprintf("say %s", pop.Value)
		}
		run(output)
		speechBuffer.Remove(pop)
	}
}

// create socket
func creatListener(addr string) net.Conn {
start:
	sock, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		goto start // pee pee poo poo lol
	}
	log.Println("dialed")
	return sock
}

// creates gsi listener
func createStateListener() *csgsi.Game {
	game := csgsi.New(3)
	if game != nil {
		return game
	}
	return nil
}

// dumb wrapper for checking state structure integrity because i cant think of a better way
func stateOK(state *csgsi.State) bool {
	if state.Previously != nil && state.Round != nil {
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

func consoleParse(toParse string) {
	// TODO owos and other stupid shit
	hashIdx := strings.Index(toParse, terribleHash)
	if hashIdx > -1 {
		checkCvars(strings.Split(toParse[:hashIdx], " "))
	}
	// find stupid symbol -> parse chat message
	if strings.Index(toParse, uniqueCode) > -1 {
		figureOutCommand(toParse)
	}
}

func figureOutCommand(msg string) {
	message, sender, isTeam := parseChat(msg)
	println(sender + " said " + message)
	/*
		if strings.Index(message, "owo") > -1 || strings.Index(message, "uwu") > -1 {
			say(getOwo(), isTeam)
		} // commented because i dont trust my own code
	*/
	if findOccurrence(message, "owo", "uwu") {
		say(getOwo(), isTeam)
	}
}

func parseChat(msg string) (message string, sender string, teamchat bool) {
	isTeam := false
	var caller string
	var output string
	codeIdx := strings.Index(msg, uniqueCode)
	if codeIdx > -1 {
		if (strings.Index(msg, "Terrorist)")) > -1 {
			isTeam = true
		}
		replacer := strings.NewReplacer(uniqueCode, "", "*DEAD*", "", "(Terrorist) ", "", "(Counter-Terrorist) ", "")
		replaced := replacer.Replace(msg)

		locationMarker := strings.LastIndex(replaced, "@") // fix
		if locationMarker > -1 {
			caller = strings.TrimSpace(replaced[:locationMarker])
		} else {
			nameEnd := strings.Index(replaced, " : ")
			caller = strings.TrimSpace(replaced[:nameEnd])
		}
		// end of SENDER get
		// start of MSG get
		msgStart := strings.Index(replaced, " : ")
		if msgStart > -1 {
			output = strings.TrimSpace(replaced[msgStart+3:])
		} else {
			println("bad msg")
			return
		}
		if len(caller) < 1 {
			caller = "<empty>"
		}
		return output, caller, isTeam
	}
	return "", "", false
}

func findOccurrence(s string, of ...string) bool {
	for t := range of {
		if strings.Index(s, of[t]) > -1 {
			return true
		}
	}
	return false
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
