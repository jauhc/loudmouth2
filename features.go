package main

import (
	"fmt"
	"log"

	"github.com/jauhc/go-csgsi"
)

func featureRadioSpam() {
	if settings.Config.Radiospam {
		run("ohn")
	}
}

func featureSendChat() {
	if speechBuffer.Len() > 1 {
		pop := speechBuffer.Front()
		poop := fmt.Sprintf("say %s", pop.Value)
		run(poop)
		speechBuffer.Remove(pop)
	}
}

func featureClan() {
	clanList := []int{7670261, 7670266, 7670268, 7670273, 7670276, 7670621, 7670634, 7670641, 7670647}
	output := fmt.Sprintf("cl_clanid 0")

	if !settings.Config.Clanid && !settings.Config.Clanfx {
		if clanIdx > 0 {
			// run this once to clear clanid after disabling
			clanIdx = -1
			fxState = true
			run("cl_clanid 0")
		}
	}
	if !settings.Config.Clanfx && settings.Config.Clanid {
		if clanIdx+1 >= len(clanList) {
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

func featureKillAnnounce(state *csgsi.State) {
	if stateOK(state) {
		if state.Previously.Player.Match_stats.Kills < state.Player.Match_stats.Kills {
			if settings.Config.Killsradio {
				run("enemydown")
			}
			if settings.Config.Kills {
				speechBuffer.PushBack(tellKill(state))
			}
			return
		}
	}
}

func featureDeathAnnounce(state *csgsi.State) {
	if stateOK(state) {
		if state.Previously.Player.Match_stats.Deaths < state.Player.Match_stats.Deaths &&
			state.Previously.Player.Match_stats.Deaths > 0 {
			log.Println("DETH")
			if settings.Config.Deaths {
				speechBuffer.PushBack(tellDeath(state))
			}
			return
		}
	}
}
