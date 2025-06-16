package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Packet struct {
	Op   uint32
	Data []byte
}

func WritePacket(conn net.Conn, op uint32, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, op)
	_ = binary.Write(&buf, binary.LittleEndian, uint32(len(jsonData)))
	_, _ = buf.Write(jsonData)

	_, err = conn.Write(buf.Bytes())
	return err
}

func ReadPacket(conn net.Conn) {
	header := make([]byte, 8)
	_, err := conn.Read(header)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	opcode := binary.LittleEndian.Uint32(header[:4])
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

	switch opcode {
	case 0: 
		fmt.Println("Handshake successful. Connection to Discord established.")
		fmt.Printf("Received Handshake Response:\n%s\n", string(prettyPayload))

	case 1: 
		fmt.Println("Activity successfully set on Discord.")
		fmt.Printf("Activity Response:\n%s\n", string(prettyPayload))
	case 3:
		fmt.Println("Heartbeat response received. Connection is alive.")
		fmt.Printf("Ping Response:\n%s\n", string(prettyPayload))
	default:
		fmt.Printf("Received unexpected response (OP %d):\n%s\n", opcode, string(prettyPayload))
	}
}