package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	HTTP_version = "HTTP/1.1"

	status_200_OK                    = "200 OK"
	status_201_Created               = "201 Created"
	status_404_Not_Found             = "404 Not Found"
	status_500_Internal_Server_Error = "500 Internal Server Error"

	// headers
	contentType   = "Content-Type"
	contentLength = "Content-Length"
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

	// init response
	res := newResponse()

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
				res.status = status_404_Not_Found
				break
			}
			res.headers[contentType] = "application/octet-stream"
			res.headers[contentLength] = strconv.Itoa(len(b))
			res.body = b

		case req.path == "/user-agent":
			if ua, ok := req.headers["User-Agent"]; ok {
				body := []byte(ua)
				res.headers[contentType] = "text/plain"
				res.headers[contentLength] = strconv.Itoa(len(body))
				res.body = body
			}

		case strings.HasPrefix(req.path, "/echo/"):
			body := []byte(req.path[6:])
			res.headers[contentType] = "text/plain"
			res.headers[contentLength] = strconv.Itoa(len(body))
			res.body = body

		default:
			res.status = status_404_Not_Found
		}

	case "POST":
		switch {
		case strings.HasPrefix(req.path, "/files/"):
			directory := os.Args[2]
			filename := req.path[len("/files"):]
			err := os.WriteFile(directory+filename, []byte(req.body), 0644)
			if err != nil {
				log.Printf("can't write to %s%s: %v\n", directory, filename, err)
				res.status = status_500_Internal_Server_Error
				break
			}
			res.status = status_201_Created
		}
	}

	// write response
	wBuf := res.build()

	_, err = conn.Write(wBuf)
	if err != nil {
		log.Println("Error writing response: ", err)
		return
	}
}
