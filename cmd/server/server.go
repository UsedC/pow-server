package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	mrand "math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const timeout = 30 * time.Second

var quotes = []string{
	"The only true wisdom is in knowing you know nothing. - Socrates",
	"Wisdom is not a product of schooling but of the lifelong attempt to acquire it. - Albert Einstein",
	"Never did nature say one thing and wisdom another. - Edmund Burke",
	"Science is organized knowledge. Wisdom is organized life. - Immanuel Kant",
	"Knowing yourself is the beginning of all wisdom. - Aristotle",
	"Turn your wounds into wisdom. - Oprah Winfrey",
	"Patience is the companion of wisdom. - Saint Augustine",
}

func main() {
	var address string
	var difficulty int

	flag.StringVar(&address, "address", "0.0.0.0:3333", "Address to listen on")
	flag.IntVar(&difficulty, "difficulty", 4, "Proof of Work difficulty level")

	flag.Parse()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer listener.Close()
	log.Println("Server is listening on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, difficulty)
	}
}

func handleConnection(conn net.Conn, difficulty int) {
	defer conn.Close()

	challenge, err := generateChallenge()
	if err != nil {
		log.Println("Error generating challenge:", err)
		return
	}

	message := fmt.Sprintf("%s:%d", challenge, difficulty)
	fmt.Fprintln(conn, message)

	conn.SetReadDeadline(time.Now().Add(timeout))
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("Error reading from client:", err)
		return
	}
	response = strings.TrimSpace(response)

	if verifyResponse(challenge, response, difficulty) {
		quote := getRandomQuote()
		fmt.Fprintln(conn, quote)
	} else {
		fmt.Fprintln(conn, "Invalid proof of work. Connection closing.")
	}
}

func generateChallenge() (string, error) {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	challenge := hex.EncodeToString(b)
	return challenge, nil
}

func verifyResponse(challenge, nonceStr string, difficulty int) bool {
	nonce, err := strconv.Atoi(nonceStr)
	if err != nil {
		log.Println("Invalid nonce format:", nonceStr)
		return false
	}
	data := challenge + nonceStr
	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:]) // converting [32]byte to []byte

	prefix := strings.Repeat("0", difficulty)
	isValid := strings.HasPrefix(hashStr, prefix)
	log.Printf("Verifying PoW: challenge=%s, nonce=%d, hash=%s, valid=%t\n", challenge, nonce, hashStr, isValid)
	return isValid
}

func getRandomQuote() string {
	mrand.Seed(time.Now().UnixNano())
	return quotes[mrand.Intn(len(quotes))]
}
