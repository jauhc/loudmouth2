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
	State      bool `json:"state"`
	Clanid     bool `json:"clanid"`
	Clanfx     bool `json:"clanfx"`
	Owo        bool `json:"owo"`
	Kills      bool `json:"kills"`
	Killsradio bool `json:"killsradio"`
	Deaths     bool `json:"deaths"`
	Greets     bool `json:"greets"`
	Ammowarn   bool `json:"ammowarn"`
	Radiospam  bool `json:"radiospam"`
}

/*
	READ THIS IF YOU WANT TO ADD MORE SETTINGS
	IF ITS A SIMPLE TOGGLE ON / OFF FEATURE, ITS EASY (KIND OF)
+---------------------------------------------------------------+
|	1. add it to struct that reads json	(LoudConfig)			|
|	2. add to checkCvars() function to make it toggleable		| <-- needs to be scalable
+---------------------------------------------------------------+
TODO reduce steps
	probably best to have struct as reference, then add missing ones to json
*/

func checkCvars(data []string) {
	fmt.Println(data)
	data[1] = strings.ToUpper(data[1])
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
		settings.Config.Killsradio = set
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
func readConfig(file string) (settings LoudSettings) {
	settingsFile, err := os.Open(file)
	ec(err)
	defer settingsFile.Close()
	byteVal, err := ioutil.ReadAll(settingsFile)
	ec(err)
	json.Unmarshal(byteVal, &settings)
	return
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
	run(fmt.Sprintf("alias loud  \"echo set 0 LIST %s\"", terribleHash))
	run("alias ohn getout")

	// parses the struct to text then gets names and throws them as aliases etc
	raw := fmt.Sprintf("%+v", settings.Config)
	col := strings.ReplaceAll(raw, " ", ":")
	col = removeAllOf(col, "{", "}")
	list := strings.Split(col, ":")
	for index := 0; index < len(list); index += 2 {
		//println(fmt.Sprintf("%v = %v", list[index], list[index+1]))
		run(fmt.Sprintf("setinfo loud_%s_o \"\"", strings.ToLower(list[index])))
		run(fmt.Sprintf("alias loud_%s_off \"echo set 0 %s %s\"",
			strings.ToLower(list[index]), strings.ToUpper(list[index]), terribleHash))
		run(fmt.Sprintf("alias loud_%s_on \"echo set 1 %s %s\"",
			strings.ToLower(list[index]), strings.ToUpper(list[index]), terribleHash))
		sleepn(30, 15)
	}

	run("echo Commands created!")
	sleep(50)
	run("loud")
}
