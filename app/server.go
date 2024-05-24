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

	status_200_OK                    = "200 OK"
	status_201_Created               = "201 Created"
	status_404_Not_Found             = "404 Not Found"
	status_500_Internal_Server_Error = "500 Internal Server Error"
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

	// read request into buffer
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading HTTP request: ", err)
		// TODO handle: respond with 500 error
		return
	}

	// parse request
	req := parseRequest(buf[:n])

	// handle URL target
	status := status_200_OK
	bodyStr := ""
	bodyByt := make([]byte, 0)

	switch req.method {
	case "GET":
		switch {
		case req.path == "/":
			break

		case strings.HasPrefix(req.path, "/files/"):
			directory := os.Args[2]
			filename := req.path[len("/files"):]
			b, err := os.ReadFile(directory + filename)
			if err != nil {
				log.Printf("can't read %s%s: %v\n", directory, filename, err)
				status = status_404_Not_Found
				break
			}
			bodyByt = b

		case req.path == "/user-agent":
			if ua, ok := req.headers["User-Agent"]; ok {
				bodyStr = ua
			}

		case strings.HasPrefix(req.path, "/echo/"):
			bodyStr = req.path[6:]

		default:
			status = status_404_Not_Found
		}
	case "POST":
		switch {
		case strings.HasPrefix(req.path, "/files/"):
			directory := os.Args[2]
			filename := req.path[len("/files"):]
			err := os.WriteFile(directory+filename, []byte(req.body), 0644)
			if err != nil {
				log.Printf("can't write to %s%s: %v\n", directory, filename, err)
				status = status_500_Internal_Server_Error
				break
			}
			status = status_201_Created
		}

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
