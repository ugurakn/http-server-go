package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
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

		handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// parse request and get URL target
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("error reading HTTP request: ", err)
		return
	}

	// extract URL path from request
	req := string(buf[:n])
	path := strings.Split(req, " ")[1]

	fmt.Println(path)

	// handle URL target
	status := "404 Not Found"
	body := ""
	switch {
	case path == "/":
		status = "200 OK"
	case strings.HasPrefix(path, "/echo/"):
		body = path[6:]
		status = "200 OK"
	}

	// write response
	buf = []byte(fmt.Sprintf("HTTP/1.1 %s\r\n\r\n", status))

	if body != "" {
		buf = formatPlain(body)
	}

	_, err = conn.Write(buf)
	if err != nil {
		log.Println("Error writing to connection: ", err)
		return
	}
}

// writePlain writes an HTTP response with the plain text s as body.
func formatPlain(body string) []byte {
	statusLine := "HTTP/1.1 200 OK\r\n"
	headers := fmt.Sprintf(
		"Content-Type: text/plain\r\nContent-Length: %v\r\n\r\n",
		len([]byte(body)),
	)
	return []byte(statusLine + headers + body)
}
