package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "localhost:4221")
	fmt.Println(l.Addr())
	if err != nil {
		log.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		// this blocks until a connection is established
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}

		// send back 200 OK
		handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// parse request and get url target
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Println("error reading HTTP request: ", err)
	}

	reqRaw := string(buf)
	path := strings.Split(reqRaw, " ")[1]

	status := "404 Not Found"
	if path == "/" {
		status = "200 OK"
	}

	buf = []byte(fmt.Sprintf("HTTP/1.1 %s\r\n\r\n", status))
	n, err := conn.Write(buf)
	if err != nil {
		log.Println("Error writing to connection: ", err)
		return
	}
	log.Printf("wrote %v bytes", n)
}
