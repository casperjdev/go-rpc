# go-rpc

A lightweight, configurable Discord Rich Presence client written in Go.

## Features

- ✅ Fully configurable via `config.json`
- ✅ Automatically refreshes presence at runtime
- ✅ Custom data embedding via BASH scripts
- ✅ No OAuth required

---

## Requirements

- Go 1.18+
- Discord desktop client (must be running)
- Unix-based OS (for default IPC path)

_Note: Windows support is not quite there yet, but it is planned._

---

## Installation

### Building from source

```bash
git clone https://github.com/casperjdev/go-rpc
cd go-rpc
go build -o go-rpc .
```

---

## Usage

### 1. Setup

Create an application on the [Discord Developer Portal](https://discord.com/developers/applications)

Save the application ID aswell as set the name of the application. This will be the name of the activity on Discord.

All of the images you will be using should also be saved there.

### 2. Configuration

Create a `config.json` file in the same directory as the build output.

The schema is adjacent to the one used by Discord's Game SDK. Same constraints apply:

```json
{
	"credentials": {
		"application_id": "123456789012345678" // Get from Discord Developer Portal
	},
	// the "activity" object follows the Discord Game SDK schema
	"activity": {
		"state": "Watching memes", // Optional (but recommended)
		"details": "Browsing Reddit", // Optional (but recommended)
		"timestamps": {
			"start": 1718454373, // Optional (epoch seconds)
			"end": 1718457973 // Optional (epoch seconds)
		},
		"assets": {
			"large_image": "big_image_key", // Optional (image key uploaded in Dev Portal)
			"large_text": "Hover text for large image", // Optional
			"small_image": "small_image_key", // Optional
			"small_text": "Hover text for small image" // Optional
		},
		"party": {
			"id": "party123", // Optional
			"size": [1, 5] // Optional (current size, max size)
		},
		"secrets": {
			"match": "match123", // Optional (for invites)
			"join": "joinSecret123", // Optional
			"spectate": "spectateSecret123" // Optional
		},
		"buttons": [
			{ "label": "My Website", "url": "https://example.com" },
			{ "label": "Join Discord", "url": "https://discord.gg/abcdef" }
		],
		"instance": true // Optional (always `true` for Rich Presence)
	},
	"constants": {
		"static": {
			"{kernel}": "uname -sr" // Either pass bash command
		},
		"dynamic": {
			"{uptime}": "./scripts/uptime.sh" // or path to runnable script
		}
	}
}
```

### 3. Embedding custom data

**Method 1:** Pass bash commands

`"{kernel}": "uname -sr"`

**Method 2:** Pass path to runnable script

`./scripts/uptime.sh`

**Static Constants** are fetched once during presence initialization and never again (data that never changes).

**Dynamic Constants** are fetched whenever the presence updates (for things like time, etc.).

Notes:

- You can name keys however you want (reccomended to use curly braces)
- Scripts _must_ be executable (`chmod +x`) and they need to output exactly one line to `stdout`
- It is reccomended to not use performance intensive scripts as they may be re-ran on every presence update.

---

## Running

Run the built executable with Discord open

> You need to be using the desktop client for Discord. Web Discord has no way to communicate through IPC.

### Running inside of a systemd service

If you want to for example automatically run this app when discord launches on your device you might want to use a **systemd service**.

Discord runs per-user, not globally. So your service must run as the same user who is logged into the graphical session (the one running Discord).

Use a user-level service, not a system-wide one.

> A reccomended setup is to have a service run a watcher script that launches the app based on Discord being active
>
> `.../.../watcher.sh`
>
> ```bash
> #!/bin/bash
> APP_NAME="go-rpc"
> APP_PATH="/path/to/app/go-rpc"
>
>
> while true; do
>     if pgrep -x Discord > /dev/null; then
>         if ! pgrep -x $(basename "$GOAPP_CMD") > /dev/null; then
>             echo "Starting Go app..."
>             "$APP_PATH" &
>         fi
>     else
>         if pgrep -x "$APP_NAME" > /dev/null; then
>             echo "Stopping Go app..."
>             pkill -x "$APP_NAME"
>         fi
>     fi
>     sleep 5
> done
> ```
>
> `.../systemd/user/go-rpc.service`
>
> ```service
> [Unit]
> Description=Run Go app when Discord is running
>
> [Service]
> ExecStart=.../.../watcher.sh
> Restart=always
>
> [Install]
> WantedBy=default.target
> ```
>
> Make sure the watcher is runnable: `chmod +x .../.../discord-watcher.sh`

---

## Editing the source code

If you wish to fine-tune certain things like update and heartbeat intervals, the code is up for grabs.

Tips:

- Make sure to not crank up the update rate too high or else your client might get flagged and blocked.
- if you prefer to write Go code instead of BASH code to write scripts for constants, you can do so by extending the `getConstants()` function in `constants.go` and the re-building your app
