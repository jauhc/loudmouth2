package main

/*	cheese.go
	contains chat related message strings bla bla bla yappp

*/

import (
	"fmt"
	"math/rand"

	"github.com/jauhc/go-csgsi"
)

var (
	gOwo = []string{
		"rawr x3 nuzzles how are you",
		"pounces on you you're so warm",
		"o3o notices you have a bulge o: someone's happy ;)",
		"nuzzles your necky wecky~ murr~ hehehe",
		"rubbies your bulgy wolgy you're so big :oooo",
		"rubbies more on your bulgy wolgy it doesn't stop growing ·///·",
		"kisses you and lickies your necky daddy likies (;",
		"nuzzles wuzzles I hope daddy really likes $:",
		"wiggles butt and squirms I want to see your big daddy meat~",
		"wiggles butt I have a little itch o3o",
		"wags tail can you please get my itch~",
		"puts paws on your chest nyea~",
		"its a seven inch itch rubs your chest can you help me pwease",
		"squirms pwetty pwease sad face I need to be punished",
		"runs paws down your chest and bites lip like I need to be punished really good~",
		"paws on your bulge as I lick my lips I'm getting thirsty",
		"I can go for some milk unbuttons your pants as my eyes glow you smell so musky :v",
		"licks shaft mmmm~ so musky drools all over your cock your daddy meat",
		"I like fondles Mr. Fuzzy Balls hehe puts snout on balls and inhales deeply",
		"oh god im so hard~ licks balls punish me daddy~",
		"nyea~ squirms more and wiggles butt I love your musky goodness",
		"bites lip please punish me licks lips nyea~",
		"suckles on your tip so good licks pre of your cock salty goodness~",
		"eyes role back and goes balls deep mmmm~ moans and suckles",
	}
)

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

func getOwo() string {
	return gOwo[rand.Intn(len(gOwo))]
}

// how to spend cpu cycles 101
func checkCommands(m chatMsg) {

	if findOccurrence(m.Message, "owo", "uwu") {
		say(getOwo(), m.Teamchat)
		return
	}

	if findOccurrence(m.Message, "d20", "!roll", "!rtd") {
		out := fmt.Sprintf("d20: %d", rand.Intn(20)+1)
		say(out, m.Teamchat)
		return
	}
}
