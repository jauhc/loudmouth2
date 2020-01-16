package main

import (
	"log"
	"os"
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

func main() {
	log.Println("start")
	t := creatListener()
	go listenerLoop(t) // thread for listening to rcon
	log.Println("Console connected!")
	gsi := createStateListener()
	for state := range gsi.Channel { // blocking loop, what we want for now
		log.Println(state.Player.Name)
	}
	log.Println("end")
}
