package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	HTTP_version = "HTTP/1.1"

	status_200_OK        = "200 OK"
	status_404_Not_Found = "404 Not Found"
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

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	// parse request and get URL target
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading HTTP request: ", err)
		return
	}

	// extract URL path from request
	req := string(buf[:n])
	path := strings.Split(req, " ")[1]

	// build headers map from request (ugly :( )
	headersRaw := strings.Split(strings.SplitN(req, "\r\n", 2)[1], "\r\n\r\n")[0]

	headers := make(map[string]string)
	for _, h := range strings.Split(headersRaw, "\r\n") {
		k, v, _ := strings.Cut(h, ": ")
		headers[k] = v
	}

	// handle URL target
	status := status_200_OK
	bodyStr := ""
	bodyByt := make([]byte, 0)
	switch {
	case path == "/":
		break

	case strings.HasPrefix(path, "/files/"):
		directory := os.Args[2]
		filename := path[len("/files"):]
		b, err := os.ReadFile(directory + filename)
		if err != nil {
			log.Printf("can't read %s%s: %v\n", directory, filename, err)
			status = status_404_Not_Found
			break
		}
		bodyByt = b

	case path == "/user-agent":
		if ua, ok := headers["User-Agent"]; ok {
			bodyStr = ua
		}

	case strings.HasPrefix(path, "/echo/"):
		bodyStr = path[6:]

	default:
		status = status_404_Not_Found
	}

	// write response
	buf = []byte(fmt.Sprintf("%s %s\r\n\r\n", HTTP_version, status))

	if bodyStr != "" {
		buf = formatPlain(bodyStr)
	} else if len(bodyByt) != 0 {
		buf = formatOctet(bodyByt)
	}

	_, err = conn.Write(buf)
	if err != nil {
		log.Println("Error writing response: ", err)
		return
	}
}

// writePlain returns a byte buffer with body as plain text.
func formatPlain(body string) []byte {
	statusLine := fmt.Sprintf("%s %s\r\n", HTTP_version, status_200_OK)
	headers := fmt.Sprintf(
		"Content-Type: text/plain\r\nContent-Length: %v\r\n\r\n",
		len(body),
	)
	return []byte(statusLine + headers + body)
}

func formatOctet(body []byte) []byte {
	statusLine := fmt.Sprintf("%s %s\r\n", HTTP_version, status_200_OK)
	headers := fmt.Sprintf(
		"Content-Type: application/octet-stream\r\nContent-Length: %v\r\n\r\n",
		len(body),
	)
	return append([]byte(statusLine+headers), body...)
}
