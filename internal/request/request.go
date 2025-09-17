package request

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
)


type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *RequestLine) Valid() bool {
	allowedMethods := []string{"POST", "GET", "PUT", "DELETE"}
	// METHOD
	isValidMethod := slices.Contains(allowedMethods, r.Method)
	log.Printf("%s", r.HttpVersion)
	isValidVersion := r.HttpVersion == "1.1"
	return isValidMethod && isValidVersion

}

const SEPERATOR = "\r\n"
var BAD_FORMAT_ERROR = fmt.Errorf("malfored request line")

func parseRequestLine(b string) (*RequestLine, error) {
	idx := strings.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, BAD_FORMAT_ERROR
	}
	startLine := b[:idx]
	parts := strings.Split(startLine, " ")

	if len(parts) != 3 {
		return nil, BAD_FORMAT_ERROR
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" {
		return nil, BAD_FORMAT_ERROR
	}

	r := &RequestLine{
		Method: parts[0],
		RequestTarget: parts[1],
		HttpVersion: httpParts[1],
	}
	if !r.Valid() {
		return nil, BAD_FORMAT_ERROR
	}
	return r, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, fmt.Errorf("unable to io.ReadAl: %w", err)
	}

	str := string(data)
	r, err := parseRequestLine(str)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *r,
	}, nil
}