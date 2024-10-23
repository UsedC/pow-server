package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const timeout = 30 * time.Second

func main() {
	var serverAddress string
	var numRequests int

	flag.StringVar(&serverAddress, "server", "0.0.0.0:3333", "Server address")
	flag.IntVar(&numRequests, "n", 1, "Number of requests to send to the server")

	flag.Parse()

	for i := 0; i < numRequests; i++ {
		log.Printf("Starting request %d of %d", i+1, numRequests)
		err := sendRequest(serverAddress)
		if err != nil {
			log.Printf("Request %d failed: %v", i+1, err)
		}
	}

	log.Println("All requests completed")
}

func sendRequest(serverAddress string) error {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return fmt.Errorf("error connecting to server: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(timeout))

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading challenge from server: %v", err)
	}
	message = strings.TrimSpace(message)

	parts := strings.Split(message, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid challenge format received from server")
	}

	challenge := parts[0]
	difficulty, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid difficulty received from server")
	}
	log.Printf("Received challenge: %s, difficulty: %d", challenge, difficulty)

	nonce := computeNonce(challenge, difficulty)
	if nonce == -1 {
		return fmt.Errorf("failed to find a valid nonce")
	}
	log.Printf("Found valid nonce: %d", nonce)

	fmt.Fprintln(conn, nonce)

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading response from server: %v", err)
	}
	log.Printf("Server response: %s", strings.TrimSpace(response))

	return nil
}

func computeNonce(challenge string, difficulty int) int {
	nonce := 0
	prefix := strings.Repeat("0", difficulty)
	startTime := time.Now()
	for {
		data := challenge + strconv.Itoa(nonce)
		hash := sha256.Sum256([]byte(data))
		hashStr := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hashStr, prefix) {
			log.Printf("PoW solved: nonce=%d, hash=%s\n", nonce, hashStr)
			return nonce
		}
		nonce++

		if time.Since(startTime) > timeout {
			return -1
		}
	}
}
