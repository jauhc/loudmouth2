# loudmouth2

rewrite of a personal project -
first iteration was c#, but now i wanted to make a native and golang seemed optimal
```
explaining files and their contents

┌───[root]
├── main.go ────────> contains the main logic for running this
├── utils.go ───────> helper functions and other boring stuff
├── settings.go ────> manages config / settings
├── cheese.go ──────> the cheesy chat message functions are here
├── features.go ────> where i try to cram most useful stuff
├── loud.json ──────> JSON formatted settings to load
├┬─── [csgo/cfg/]
│└─── gamestate_integration_loudmouth2.cfg ──> tells the game connection data
```

## "i somehow ended up here, what is this?"
in short: a csgo spambot

### longer:
csgo (source engine in general, at some point after l4d) allows one to connect to its console remotely by setting `-netconport` launch parameter value.

in raw form its essentially just opening a connection to said port then listen for data (requires authentication with `PASS passwordhere` if using `-netconpassword` (not sure if its limited to local connections))

### known issues
- game crashes on linux upon connection