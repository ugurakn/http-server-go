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

	// build headers map from request
	headersRaw := strings.Split(strings.SplitN(req, "\r\n", 2)[1], "\r\n\r\n")[0]
	headers := make(map[string]string)
	for _, h := range strings.Split(headersRaw, "\r\n") {
		k, v, _ := strings.Cut(h, ": ")
		headers[k] = v
	}

	// for k, v := range headers {
	// 	fmt.Printf("%s:%s.\n", k, v)
	// }

	// handle URL target
	status := "200 OK"
	body := ""
	switch {
	case path == "/":
		break
	case path == "/user-agent":
		if ua, ok := headers["User-Agent"]; ok {
			body = ua
		}
	case strings.HasPrefix(path, "/echo/"):
		body = path[6:]
	default:
		status = "404 Not Found"
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

// writePlain writes a byte buffer with body as plain text.
func formatPlain(body string) []byte {
	statusLine := "HTTP/1.1 200 OK\r\n"
	headers := fmt.Sprintf(
		"Content-Type: text/plain\r\nContent-Length: %v\r\n\r\n",
		len([]byte(body)),
	)
	return []byte(statusLine + headers + body)
}
