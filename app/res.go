package main

import (
	"fmt"
)

type response struct {
	status  string
	headers map[string]string
	body    []byte
}

func newResponse() *response {
	return &response{
		status:  status_200_OK,
		headers: make(map[string]string),
	}
}

func (res *response) build() []byte {
	sLine := res.statusLine()
	headers := res.getHeaders()

	totalLen := len(sLine) + len(headers) + len(res.body)
	r := make([]byte, totalLen)

	var i int
	i += copy(r[i:], sLine)
	i += copy(r[i:], headers)
	i += copy(r[i:], res.body)

	return r
}

func (res *response) statusLine() []byte {
	//<http_version> <status> CRLF
	s := fmt.Sprintf("%s %s\r\n", HTTP_version, res.status)
	return []byte(s)
}

func (res *response) getHeaders() []byte {
	// each header:
	//<headerKey>: <headerVal>CRLF
	// extra CRLF at the end
	var s string
	for k, v := range res.headers {
		s += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	s += "\r\n"

	return []byte(s)
}
