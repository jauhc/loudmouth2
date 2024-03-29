package main

/*	settings.go
	for handling config/settings so it doesnt clutter utils.go
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/andygrunwald/vdf"
)

// LoudSettings represents root level config: user, configs
type LoudSettings struct {
	User    string     `json:"user"`    // your steamcommunity id
	Gsiport string     `json:"gsiport"` // port used by GSI
	Netport string     `json:"netport"` // console port
	Pass    string     `json:"pass"`    // console connection password, same as -netconpassword
	Config  LoudConfig `json:"config"`
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
*/

// TODO
// https://github.com/andygrunwald/vdf
func getTelnetParams() {
	pa, err := getRegistryValue("SOFTWARE\\Valve\\Steam", "SteamPath")
	ec(err)
	file := fmt.Sprintf("[%s/userdata/%s/config/localconfig.vdf]", pa, getSteamID(false))
	f, err := os.Open(file)
	ec(err)

	p := vdf.NewParser(f)
	m, err := p.Parse()
	ec(err)

	// holy SHIT
	rawparams := fmt.Sprintf("%s", m["UserLocalConfigStore"].(map[string]interface{})["Software"].(map[string]interface{})["Valve"].(map[string]interface{})["Steam"].(map[string]interface{})["apps"].(map[string]interface{})["730"].(map[string]interface{})["LaunchOptions"])
	params := strings.Split(rawparams, " ")

	for i := 0; i < len(params); i++ {
		if params[i] == "-netconport" {
			settings.Netport = params[i+1]
			println("Port found in launch params:", settings.Netport)
		}
		if params[i] == "-netconpassword" {
			settings.Pass = params[i+1]
			println("PASS found in params:", settings.Pass)
		}
	}
	saveConfig()
}

func checkCvars(data []string) {
	if len(data) < 2 {
		return
	}
	// clean strings
	for i := 0; i < len(data); i++ {
		data[i] = removeAllOf(data[i], " ")
	}
	data[1] = strings.ToUpper(data[1])
	set := false
	if data[0] == "1" {
		set = true
	}
	switch data[1] {
	case "LIST":
		raw := fmt.Sprintf("%+v", settings.Config)
		raw = removeAllOf(raw, "{", "}")
		list := strings.Split(raw, " ")
		// startFancy()
		// defer endFancy()
		startColour()
		for index := 0; index < len(list); index++ {
			/*	make string to `setting:bool`
				split with : to array
				eval val[1]	*/
			ss := strings.Split(removeAllOf(list[index], " "), ":")
			b, err := strconv.ParseBool(ss[1])
			ec(err)
			if b {
				useColour("2fb54aFF")
				run(fmt.Sprintf("echo %s\n", strings.ToLower(ss[0])))
			} else {
				useColour("d15532FF")
				run(fmt.Sprintf("echo %s\n", strings.ToLower(ss[0])))
			}

			/// JESUS CHRIST I NEED TO RETHINK THIS
			// can spam echo since its client command, server doesnt care
			// split with : -> 2nd value boolean
			// run(fmt.Sprintf("echo %v\n", strings.ToLower(list[index])))
		}
		endColour()
		break

	case "STATE":
		settings.Config.State = set
		break

	case "AMMOWARN":
		settings.Config.Ammowarn = set
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

	case "DMGREPORT": // prob never done since cant check for alive players
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

	run(fmt.Sprintf("con_filter_text_out %s", terribleHash))
	run("echo Commands created!")
	sleep(50)
	run("loud")
}
