package main

import (
	"fmt"
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
	conn, err := Connect()
	if err != nil {
		fmt.Println("Failed to Connect to Discord IPC")
		os.Exit(1)
	}

	// Handshake
	applicationID := config.Credentials.ApplicationID
	if applicationID == "" {
		fmt.Println("Missing APPLICATION_ID environment variable")
		os.Exit(1)
	}

	handshake := map[string]any{
		"v":         1,
		"client_id": applicationID,
	}

	if err := SendPacket(conn, 0, handshake); err != nil {
		fmt.Println("Failed to send handshake:", err)
		return
	}
	ReadPacket(conn, "Handshake Output: ");

	// Update loop
	update := time.NewTicker(15 * time.Second)
	defer update.Stop()

	go func() {
		// Fetch static constants once
		staticConstants := getConstants(config.Constants.Static)

		// Initialize presence once
		activityClone, err := cloneActivity(config.Activity)
		if err != nil {
			fmt.Println("Failed to clone activity: ", err)
		} else {
			// Fetch dynamic constants and merge with static
			constants := MergeConstants(staticConstants, getConstants(config.Constants.Dynamic))

			ProcessConstants(&activityClone, constants)

			activityPayload := map[string]any{
				"cmd": "SET_ACTIVITY",
				"args": map[string]any{
					"pid":      os.Getpid(),
					"activity": activityClone,
				},
				"nonce": fmt.Sprintf("%d", time.Now().UnixNano()),
			}
			
			if err := SendPacket(conn, 1, activityPayload); err != nil {
				fmt.Println("Failed to update presence:", err)
			} else {
				fmt.Println("Presence updated at", time.Now().Format(time.RFC1123))
			}
			ReadPacket(conn, "Resulting Presence Object: ")
		}

		// ...Then run every ticker cycle
		for range update.C {
			activityClone, err := cloneActivity(config.Activity)
			if err != nil {
				fmt.Println("Failed to clone activity: ", err)
			} else {
				// Fetch dynamic constants and merge with static
				constants := MergeConstants(staticConstants, getConstants(config.Constants.Dynamic))
				ProcessConstants(&activityClone, constants)

				activityPayload := map[string]any{
					"cmd": "SET_ACTIVITY",
					"args": map[string]any{
						"pid": os.Getpid(),
						"activity": activityClone,
						},
					"nonce": fmt.Sprintf("%d", time.Now().UnixNano()),
				}

				if err := SendPacket(conn, 1, activityPayload); err != nil {
					fmt.Println("Failed to update presence:", err)
				} else {
					fmt.Println("Presence updated at", time.Now().Format(time.RFC1123))
				}
			}
			
		}
	}()

	// Heartbeat loop
	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	for range heartbeat.C {
		if err := SendPacket(conn, 3, map[string]any{}); err != nil {
			fmt.Println("Heartbeat failed:", err)
			return
		}
	}
}
