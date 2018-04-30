package main

import (
	"fmt"
	"net"
	"encoding/binary"
	"bytes"
)

const PacketSize = 65535

// Packet Type RFC 1350 section 5
const (
	RRQ   = iota + 1 // Read Request
	WRQ              // Write Request
	DATA             // Data
	ACK              // Acknowledgment
	ERROR            // Error
)

// Here's the plan....
// Listen to port 3333 and form a map of blocks which will form an entire file
// Once we've received the file, display it on the terminal.
func main() {
	var file map[uint16][]byte

	// Get information about the local host and the port
	serverAddress, err := net.ResolveUDPAddr("udp", ":3333")
	if err != nil {
		fmt.Println("Error getting our address: ", err)
		return
	}

	// Open the port to listen!
	conn, err := net.ListenUDP("udp", serverAddress)
	if err != nil {
		panic(err)
	}

	fileInProgress := false
	// Loop forever reading from the port
	for {
		buffer := make([]byte, PacketSize)

		n, remoteAddress, err := conn.ReadFromUDP(buffer[:])

		if err != nil {
			fmt.Println("Error reading from network: ", err)
			return
		}

		// Pass the info we received to a helper to work out the details
		ok, blockNumber, filename, _, data := processPacket(conn, remoteAddress, buffer[:n], n)

		// We processed the packet but it's nothing we care about
		if !ok {
			continue
		}

		// We're not tracking an incoming file but we got some data that is mid file
		if !fileInProgress && blockNumber > 0 {
			fmt.Println("Not expecting mid file data")
			continue
		}

		// We're not tracking a file and we got a starting block
		if !fileInProgress && blockNumber == 0 {
			fmt.Println("Starting new file ", filename)
			file = make(map[uint16][]byte)
			fileInProgress = true
			continue
		}

		// Received a block that we are happy to get
		if fileInProgress {
			file[blockNumber] = data

			// Got the final block
			if len(data) < 512 {
				fileInProgress = false
				printFile(file, blockNumber)
			}
		}
	}
}

// processPracket takes the buffer from the reading the port and processes it, sending an appropriate reply to the
// sender.
// Returns a success flag (ok) if we got something we care about (Write request and data).  Other return data includes
// file names, block number and data blob
func processPacket(conn net.PacketConn, remoteAddress net.Addr, buffer []byte, n int) (ok bool, blockNumber uint16, filename string, mode string, data []byte) {
	ok = false

	// Did we get enough data for an Opcode?
	if n < 2 {
		return
	}

	opCode := binary.BigEndian.Uint16(buffer[0:2])
	switch opCode {
	case RRQ:
		// Not supporting read functions
		sendError(conn, remoteAddress, "SORRY! Reading not supported")

	case WRQ:
		// Starting a write request.
		content := bytes.Split(buffer[2:], []byte{0})
		filename = string(content[0])
		mode = string(content[1])

		sendAck(conn, remoteAddress, 0)
		ok = true

	case DATA:
		// Received more data
		blockNumber = uint16(binary.BigEndian.Uint16(buffer[2:4]))
		data = buffer[4:]
		sendAck(conn, remoteAddress, blockNumber)
		ok = true

	case ACK:
		// Do nothing for an ACK
		// blockNumber := binary.BigEndian.Uint16(buffer[2:4])

	case ERROR:
		// We received an error from the client
		errorCode := binary.BigEndian.Uint16(buffer[2:4])
		errorMessage := string(buffer[4:])
		fmt.Println("Error", errorCode, errorMessage)

	default:
		return
	}

	return
}

// toBytes takes the input uint16 array and converts it to an array of bytes
func toBytes(input []uint16) []byte {
	output := make([]byte, len(input)*2)

	for index, element := range input[:] {
		outputIndex := index * 2
		binary.BigEndian.PutUint16(output[outputIndex:outputIndex+2], uint16(element))
	}
	return output
}

// sendError sends the error op code "undefined" with an explanatory  message to the specified remote address via
// the specified connection
func sendError(conn net.PacketConn, remoteAddress net.Addr, message string) {
	errorMessage := []byte(message)
	errorCode := toBytes([]uint16{ERROR, 0})
	errorResponse := append(errorCode, errorMessage...)

	_, err := conn.WriteTo(append(errorResponse, 0), remoteAddress)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// sendAck sends the ACK op code and a block number to the specified remote address via the specified connection
func sendAck(conn net.PacketConn, remoteAddress net.Addr, blockNumber uint16) {
	_, err := conn.WriteTo(toBytes([]uint16{ACK, uint16(blockNumber)}), remoteAddress)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// printFile takes the map that was created in the main loop and displays it to the console.
// if there is a block missing it displays an error and returns false.  If all blocks are there, it returns true
func printFile(file map[uint16][]byte, lastBlock uint16) bool {

	for i := uint16(1); i <= lastBlock; i++ {
		if data, gotBlock := file[i]; gotBlock {
			fmt.Print(string(data))
		} else {
			fmt.Println("ERROR! Didn't get all the blocks!")
			return false
		}
	}
	return true
}



