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
	conColour    = "00000000"
)

func createTimers() {
	go gayTicker()
}

// dum sleep wrapper because im lazy
func sleep(ms int) {
	if ms > 0 {
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

// sleep with some noise
func sleepn(ms int, noise int) {
	n := rand.Intn(noise)
	time.Sleep(time.Duration(ms+n) * time.Millisecond)
}

// short for error check
func ec(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func startColour() {
	run("con_filter_enable 1; con_filter_text_out Setting")
	sleep(45) // somehow necessary...
}

func useColour(s string) {
	if conColour == s {
		return
	}
	conColour = s
	run(fmt.Sprintf("log_color console %s", conColour)) // 00D900FF
}

func endColour() {
	run(fmt.Sprintf("log_color console 00000000; con_filter_text_out %s", terribleHash))
}

func doFancy(s string) {
	startColour()
	useColour("FF00FFFF")
	run(fmt.Sprintf("echo %s", s))
	endColour()
}

// removes all give substrs of s
func removeAllOf(s string, r ...string) (t string) {
	t = strings.ReplaceAll(s, r[0], "")
	for i := 1; i < len(r); i++ {
		t = strings.ReplaceAll(t, r[i], "")
	}
	return
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
	message, sender, location, isTeam := parseChat(msg)
	if len(location) > 1 {
		print("[@" + location + "] ")
	}
	println("(" + sender + "): " + message)
	/*
		if strings.Index(message, "owo") > -1 || strings.Index(message, "uwu") > -1 {
			say(getOwo(), isTeam)
		} // commented because i dont trust my own code
	*/
	if findOccurrence(message, "owo", "uwu") {
		say(getOwo(), isTeam)
	}
}

// probably better just to struct this
func parseChat(msg string) (message string, sender string, location string, teamchat bool) {

	codeIdx := strings.Index(msg, uniqueCode)
	// var isDead = false
	var isTeam = false
	sender = ""
	location = ""

	// actual checking code begins here //
	if codeIdx > -1 { // is a player message

		// figure out name
		if strings.Index(msg, "*") == 0 { // has *DEAD* at start
			deadmarkerLen := strings.Index(msg[1:], "*")
			if deadmarkerLen > -1 {
				sender = msg[deadmarkerLen+3 : codeIdx]
			}
		} else { // does NOT have *DEAD* at start
			sender = msg[:codeIdx]
		}

		/*
			if strings.Index(msg, "*DEAD*") > -1 { // is from dead
				isDead = true
			}
		*/

		if strings.Index(msg, "T)") > -1 { // is from team
			isTeam = true
		}

		if strings.Index(msg[codeIdx:], "@") > -1 { // has location
			endofLocation := strings.Index(msg[codeIdx:], ":")
			location = msg[codeIdx+6 : codeIdx+endofLocation]
			location = location[:strings.LastIndex(location, "(")-1] // cleanup
			print("@[" + location + "] ")
		}

		/*
				if isTeam {
					print("(TEAM)")
				}
				if isDead {
					print("<DEAD>")
				}
				if !isDead && !isTeam {
					//print("[ALL]")
				}

			print("{" + sender + "}")
		*/

		// remove the shit we already have
		withoutName := msg[codeIdx+4:]
		startOfMsg := strings.Index(withoutName, ":")
		//println(withoutName[startOfMsg+2:])
		message = withoutName[startOfMsg+2:]
		return message, sender, location, isTeam

	}
	return "", "", "", false
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
