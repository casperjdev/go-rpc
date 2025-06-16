package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/viper"
)

func main() {
	// Load values from config
	var config Config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Unable to decode into struct, %v\n", err)
		os.Exit(1)
	}
	
	// Connecting to IPC
	socketPath := os.Getenv("DISCORD_IPC_PATH")
	if socketPath == "" {
		socketPath = "/run/user/1000/discord-ipc-0" // default for Unix systems
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Println("Failed to connect to Discord IPC:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Send a Handshake
	applicationID := config.Credentials.ApplicationID
	if applicationID == "" {
		fmt.Println("Missing APPLICATION_ID environment variable")
		os.Exit(1)
	}

	handshake := map[string]any{
		"v":         1,
		"client_id": applicationID,
	}

	if err := WritePacket(conn, 0, handshake); err != nil {
		fmt.Println("Failed to send handshake:", err)
		return
	}
	ReadPacket(conn);

	// Update ticker
	updateTicker := time.NewTicker(30 * time.Second)
	defer updateTicker.Stop()

	go func() {
		// Initialize presence once, then run every ticker cycle
		activityClone, err := cloneActivity(config.Activity)
		if err != nil {
			fmt.Println("Failed to clone activity: ", err)
		} else {
			ProcessConstants(&activityClone, getConstants())

			activityPayload := map[string]any{
				"cmd": "SET_ACTIVITY",
				"args": map[string]any{
					"pid":      os.Getpid(),
					"activity": activityClone,
				},
				"nonce": fmt.Sprintf("%d", time.Now().UnixNano()),
			}

			// Show resulting object for debugging purposes
			payloadBytes, err := json.MarshalIndent(activityPayload, "", "  ")
			if err != nil {
				fmt.Println("Failed to parse activity:", err)
			} else {
				fmt.Println("Prepared full activty payload:")
				fmt.Println(string(payloadBytes))
			}
			
			setActivity(conn, activityPayload)
		}

		for range updateTicker.C {
			activityClone, err := cloneActivity(config.Activity)
			if err != nil {
				fmt.Println("Failed to clone activity: ", err)
			}
			ProcessConstants(&activityClone, getConstants())

			activityPayload := map[string]any{
				"cmd": "SET_ACTIVITY",
				"args": map[string]any{
					"pid": os.Getpid(),
					"activity": activityClone,
					},
				"nonce": fmt.Sprintf("%d", time.Now().UnixNano()),
			}

			setActivity(conn, activityPayload)
		}
	}()



	// Step 3: Heartbeats
	fmt.Println("Rich Presence set. Sending heartbeats...")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := WritePacket(conn, 3, map[string]any{}); err != nil {
			fmt.Println("Heartbeat failed:", err)
			return
		}
	}
}
