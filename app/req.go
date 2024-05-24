package main

import (
	"fmt"
	"strings"
)

type request struct {
	method  string
	path    string
	headers map[string]string
	body    string
}

// parseRequest accepts a request buffer and
// builds and returns a pointer to a request instance.
func parseRequest(buf []byte) *request {
	reqRaw := string(buf)
	req := new(request)

	reqLine := strings.Split(reqRaw, "\r\n")[0]

	req.method = strings.Split(reqLine, " ")[0]
	req.path = strings.Split(reqLine, " ")[1]

	// build headers map
	headersRaw := strings.Split(
		strings.SplitN(reqRaw, "\r\n", 2)[1],
		"\r\n\r\n",
	)[0]

	req.headers = make(map[string]string)
	for _, h := range strings.Split(headersRaw, "\r\n") {
		k, v, _ := strings.Cut(h, ": ")
		req.headers[k] = v
	}

	// body
	req.body = strings.Split(reqRaw, "\r\n\r\n")[1]

	fmt.Println(req.method, req.path)
	fmt.Println(req.headers)
	fmt.Println(req.body)

	return req
}
