# loudmouth2

rewrite of a personal project
first iteration was c#, but now i wanted to make a native and golang seemed optimal
```
explaining files and their contents

┌──[root]
├── main.go ────────> contains the main logic for running this
├── utils.go ───────> helper functions and other boring stuff
├── settings.go ────> manages config / settings
├── cheese.go ──────> the cheesy chat message functions are here
├── features.go ────> where i try to cram most useful stuff
├── loud.json ──────> JSON formatted settings to load
├┬── [csgo/cfg/]
│└── gamestate_integration_loudmouth2.cfg ──> tells the game connection data
```

## "i somehow ended up here, what is this?"
in short: a csgo spambot

longer: see [main repo](https://github.com/jauhc/loud)