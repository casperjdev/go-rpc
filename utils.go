package main

import (
	"fmt"
	"net"
	"time"

	"encoding/json"
)

func setActivity(conn net.Conn, payload map[string]any ) {
	err := WritePacket(conn, 1, payload) // assuming this is your combined version
	
	if err != nil {
		fmt.Println("Failed to update presence:", err)
	} else {
		fmt.Println("Presence updated at", time.Now().Format(time.RFC1123))
	}
}

func cloneActivity(activity DiscordActivity) (DiscordActivity, error) {
	var clone DiscordActivity

	data, err := json.Marshal(activity)
	if err != nil {
		return clone, err
	}

	err = json.Unmarshal(data, &clone)
	return clone, err
}


