package main

/*	settings.go
	for handling config/settings so it doesnt clutter utils.go
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// LoudSettings represents root level config: user, configs
type LoudSettings struct {
	User   string     `json:"user"` // your steamcommunity id
	Pass   string     `json:"pass"` // console connection password, same as -netconpassword
	Config LoudConfig `json:"config"`
}

// password is static until i can figure out how to automate this
// shit without forcing this tool to be a csgo launcher

// LoudConfig contains settings hardcoded until UI support?
type LoudConfig struct {
	State     bool `json:"state"`
	Clanid    bool `json:"clanid"`
	Clanfx    bool `json:"clanfx"`
	Owo       bool `json:"owo"`
	Kills     bool `json:"kills"`
	Deaths    bool `json:"deaths"`
	Greets    bool `json:"greets"`
	Ammowarn  bool `json:"ammowarn"`
	Radiospam bool `json:"radiospam"`
}

/*
	READ THIS IF YOU WANT TO ADD MORE SETTINGS
	IF ITS A SIMPLE TOGGLE ON / OFF FEATURE, ITS EASY (KIND OF)
+---------------------------------------------------------------+
|	1. add it to .json file following previous examples			|
|	2. add it to struct that reads json	(LoudConfig)			|
|	3. add to checkCvars() function to make it toggleable		|
|	4. add an alias to it in createConsoleCommands()			|
+---------------------------------------------------------------+
*/

func checkCvars(data []string) {
	fmt.Println(data)
	set := false
	if data[0] == "1" {
		set = true
	}
	switch data[1] {
	case "LIST":
		raw := fmt.Sprintf("%+v", settings.Config)
		list := strings.Split(raw, " ")
		for index := 0; index < len(list); index++ {
			run(fmt.Sprintf("echo %v \n", list[index]))
		}
		break

	case "STATE":
		settings.Config.State = set
		break

	case "OWO":
		settings.Config.Owo = set
		break

	case "CLAN":
		settings.Config.Clanid = set
		run("cl_clanid 0")
		break

	case "RADIOSPAM":
		settings.Config.Radiospam = set
		break

	case "CLANFX":
		settings.Config.Clanfx = set
		break

	case "DMGREPORT":
		run("echo #unimplemented")
		break

	case "KILLS":
		settings.Config.Kills = set
		break

	case "KILLSRADIO":
		run("echo #unimplemented")
		break

	case "DETH":
		settings.Config.Deaths = set
		break

	case "GREET":
		settings.Config.Greets = set
		break

	default:
		run("echo somehow you broke the settings?")
		break
	}
	saveConfig()
}

// read and apply user, configs etc
func readConfig(file string) LoudSettings {
	var settings LoudSettings
	settingsFile, err := os.Open(file)
	ec(err)
	defer settingsFile.Close()
	byteVal, err := ioutil.ReadAll(settingsFile)
	ec(err)
	json.Unmarshal(byteVal, &settings)
	return settings
}

// save current config
func saveConfig() {
	data, err := json.MarshalIndent(settings, "", " ")
	ec(err)
	err = ioutil.WriteFile(configFile, data, 0644)
	ec(err)
}

// create console commands / aliases
func createConsoleCommands() {
	// christ the sprintf makes this awful to read but its a clever 1 liner
	run(fmt.Sprintf("alias loud  \"echo set 0 LIST %s\"", terribleHash))
	run("alias ohn getout")

	run("setinfo loud_state_o \"\"")
	run(fmt.Sprintf("alias loud_state_off \"echo set 0 STATE %s\"", terribleHash))
	run(fmt.Sprintf("alias loud_state_on \"echo set 1 STATE %s\"", terribleHash))
	sleep(50)

	run("setinfo loud_owo_o \"\"")
	run(fmt.Sprintf("alias loud_owo_off \"echo set 0 OWO %s\"", terribleHash))
	run(fmt.Sprintf("alias loud_owo_on \"echo set 1 OWO %s\"", terribleHash))
	sleep(50)

	run("setinfo loud_radiospam_o \"\"")
	run(fmt.Sprintf("alias loud_radiospam_off \"echo set 0 RADIOSPAM %s\"", terribleHash))
	run(fmt.Sprintf("alias loud_radiospam_on \"echo set 1 RADIOSPAM %s\"", terribleHash))
	sleep(50)

	run("setinfo loud_dmgreport_o \"\"")
	run(fmt.Sprintf("alias loud_dmgreport_off \"echo set 0 DMGREPORT %s\"", terribleHash))
	run(fmt.Sprintf("alias loud_dmgreport_on \"echo set 1 DMGREPORT %s\"", terribleHash))
	sleep(50)

	run("setinfo loud_kills_o \"\"")
	run(fmt.Sprintf("alias loud_kills_off \"echo set 0 KILLS %s\"", terribleHash))
	run(fmt.Sprintf("alias loud_kills_on \"echo set 1 KILLS %s\"", terribleHash))
	sleep(50)

	run("setinfo loud_killradio_o \"\"")
	run(fmt.Sprintf("alias loud_killradio_off \"echo set 0 KILLSRADIO %s\"", terribleHash))
	run(fmt.Sprintf("alias loud_killradio_on \"echo set 1 KILLSRADIO %s\"", terribleHash))
	sleep(50)

	run("setinfo loud_death_o \"\"")
	run(fmt.Sprintf("alias loud_death_off \"echo set 0 DETH %s\"", terribleHash))
	run(fmt.Sprintf("alias loud_death_on \"echo set 1 DETH %s\"", terribleHash))
	sleep(50)

	run("setinfo loud_greet_o \"\"")
	run(fmt.Sprintf("alias loud_greet_off \"echo set 0 GREET %s\"", terribleHash))
	run(fmt.Sprintf("alias loud_greet_on \"echo set 1 GREET %s\"", terribleHash))
	sleep(50)

	run("setinfo loud_clan_o \"\"")
	run(fmt.Sprintf("alias loud_clan_off \"echo set 0 CLAN %s\"", terribleHash))
	run(fmt.Sprintf("alias loud_clan_on \"echo set 1 CLAN %s\"", terribleHash))
	sleep(50)

	run("setinfo loud_clan_wave_o \"\"")
	run(fmt.Sprintf("alias loud_clan_wave_off \"echo set 0 CLANFX %s\"", terribleHash))
	run(fmt.Sprintf("alias loud_clan_wave_on \"echo set 1 CLANFX %s\"", terribleHash))
	sleep(50)

	run("echo Commands created!")
	sleep(50)
	run("loud")
}
