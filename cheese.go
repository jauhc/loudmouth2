package main

/*	cheese.go
	contains chat related message strings bla bla bla yappp
	
*/

import (
	"fmt"
	"math/rand"

	"github.com/jauhc/go-csgsi"
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
