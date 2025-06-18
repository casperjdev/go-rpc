package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
)

type Packet struct {
	Op   uint32
	Data []byte
}

func Connect() (net.Conn, error) {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	// Fallback to /run/user/UID if XDG_RUNTIME_DIR is not set
	if runtimeDir == "" {
		u, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("failed to determine current user: %w", err)
		}
		runtimeDir = filepath.Join("/run/user", u.Uid)
	}
	const maxIPC = 10

	for i := range maxIPC {
		ipcPath := filepath.Join(runtimeDir, fmt.Sprintf("discord-ipc-%d", i))

		if _, err := os.Stat(ipcPath); err != nil {
			continue // skip non-existent socket
		}

		conn, err := net.Dial("unix", ipcPath)
		if err != nil {
			continue // skip failed attempts
		}

		fmt.Printf("Connected to Discord IPC socket: %s\n", ipcPath)
		return conn, nil // success
	}

	return nil, fmt.Errorf("no available Discord IPC socket found")
}

func SendPacket(conn net.Conn, op uint32, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, op)
	_ = binary.Write(&buf, binary.LittleEndian, uint32(len(jsonData)))
	_, _ = buf.Write(jsonData)

	_, _ = conn.Write(buf.Bytes())
	return err
}

func ReadPacket(conn net.Conn, message string) {
	header := make([]byte, 8)
	_, err := conn.Read(header)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// opcode := binary.LittleEndian.Uint32(header[:4])
	length := binary.LittleEndian.Uint32(header[4:])

	payload := make([]byte, length)
	_, err = conn.Read(payload)
	if err != nil {
		fmt.Println("Error reading payload:", err)
		return
	}

	payloadStr := string(payload)

	var prettyPayload []byte
	var jsonErr error
	var jsonData interface{}
	if json.Valid([]byte(payloadStr)) {
		err := json.Unmarshal([]byte(payloadStr), &jsonData)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

		prettyPayload, jsonErr = json.MarshalIndent(jsonData, "", "  ")
		if jsonErr != nil {
			fmt.Println("Error formatting JSON:", jsonErr)
			return
		}
	}

	fmt.Println(message, string(prettyPayload));
}