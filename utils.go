package main

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"

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

	t            = creatListener(fmt.Sprintf(":%s", settings.Netport))
	settings     = readConfig(configFile)
	terribleHash = generateTerribleHash(9)
	conColour    = "00000000"
)

func createTimers() {
	go gayTicker()
}

// method to run func with ease
type funcDef func()

func do(howmany int, fun funcDef) {
	for i := 0; i < howmany; i++ {
		fun()
	}
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
				log.Println(fmt.Sprintf("aaaa %#v", state.Player.State.Flashed))
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

func getRegistryValue(path string, key string) (string, error) {

	k, err := registry.OpenKey(registry.CURRENT_USER, path, 1)
	ec(err)
	defer k.Close()
	var buf [128]byte // questionable
	_, keytype, err := k.GetValue(key, buf[:])
	ec(err)
	if keytype == 4 { // DWORD
		data, _, err := k.GetIntegerValue(key)
		ec(err)
		u := strconv.FormatUint(data, 10)
		return u, nil
	} else if keytype == 1 { // SZ
		data, _, err := k.GetStringValue(key)
		ec(err)
		return data, nil
	}
	return "", errors.New("registry read went wrong?")
	// need more types, not for now though
}

// true for communityid, false for steam3id
func getSteamID(community bool) string {
	user, err := getRegistryValue("SOFTWARE\\Valve\\Steam\\ActiveProcess", "ActiveUser")
	ec(err)
	if community == true {
		// return COMMUNITY ID
		co, err := strconv.ParseUint(user, 10, 32)
		ec(err)
		return fmt.Sprintf("7656%d\n", co+1197960265728)
	}
	// return STEAM3ID
	return user
}

// prevents shit from spectated players
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
		formatChatMessage(toParse)
	}
}

type chatMsg struct {
	Message  string
	Sender   string
	Location string
	Teamchat bool
	IsDead   bool
}

func formatChatMessage(msg string) {
	m, err := parseChat(msg)

	if err != nil {
		log.Fatalln(err)
		return
	}

	/*
		if len(m.Location) > 1 {
			print("[@" + m.Location + "] ")
		}
		println("(" + m.Sender + "): " + m.Message)
	*/

	checkCommands(m) // checks for commands
	// TODO make a more scriptlike format for this (external file?)

}

func parseChat(msg string) (m chatMsg, err error) {
	//
	// using the textmod below
	// https://gist.github.com/xPaw/056b29be7ae9c143ed623a9c4c10cf50#file-csgo_bananagaming-txt
	//
	codeIdx := strings.Index(msg, uniqueCode)
	m.Sender = ""
	m.Location = ""

	// actual checking code begins here //
	if codeIdx > -1 { // is a player message

		// figure out name
		if strings.Index(msg, "*") == 0 { // has *DEAD* at start
			deadmarkerLen := strings.Index(msg[1:], "*")
			if deadmarkerLen > -1 {
				m.Sender = msg[deadmarkerLen+3 : codeIdx]
			}
		} else { // does NOT have *DEAD* at start
			m.Sender = msg[:codeIdx]
		}

		if strings.Index(msg, "*DEAD*") > -1 { // is from dead
			m.IsDead = true
		}

		if strings.Index(msg, "T)") > -1 { // is from team
			m.Teamchat = true
		}

		if strings.Index(msg[codeIdx:], "@") > -1 { // has location
			endofLocation := strings.Index(msg[codeIdx:], ":")
			m.Location = msg[codeIdx+6 : codeIdx+endofLocation]
			m.Location = m.Location[:strings.LastIndex(m.Location, "(")-1] // cleanup
			print("@[" + m.Location + "] ")
		}

		/*
			if m.Teamchat {
				print("(TEAM)")
			}
			if m.IsDead {
				print("<DEAD>")
			}
			if !m.IsDead && !m.Teamchat {
				//print("[ALL]")
			}

			print("{" + m.Sender + "}")
		*/

		// remove the shit we already have
		withoutName := msg[codeIdx+4:]
		startOfMsg := strings.Index(withoutName, ":")
		m.Message = withoutName[startOfMsg+2:]
		return m, nil

	}
	return m, errors.New("Unable to parse chat message")
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
