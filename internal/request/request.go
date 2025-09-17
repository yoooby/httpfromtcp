package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"slices"
)


type RequestState int

const (
	RequestInitialized RequestState = 0;
	RequestDone RequestState = 1;
)

type Request struct {
	RequestLine RequestLine
	State RequestState
}

func newRequest() *Request {
	return &Request{
		State: RequestInitialized,
	}
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
	outer:
	for {
		switch r.State {
		case RequestInitialized:
			n, rl, err := parseRequestLine(data)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			read += n
			r.RequestLine = *rl

			r.State = RequestDone
		case RequestDone:
			break outer
		}
	}

	return read, nil

}

func (r *Request) isDone() bool {
	return r.State == RequestDone
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

var SEPERATOR = []byte("\r\n")
var BAD_FORMAT_ERROR = fmt.Errorf("malfored request line")

func parseRequestLine(b []byte) (int, *RequestLine, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return 0, nil, nil
	}
	startLine := b[:idx]
	parts := bytes.Split(startLine, []byte(" "))

	if len(parts) != 3 {
		return 0, nil, BAD_FORMAT_ERROR
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" {
		return 0, nil, BAD_FORMAT_ERROR
	}

	r := &RequestLine{
		Method: string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion: string(httpParts[1]),
	}
	if !r.Valid() {
		return 0, nil, BAD_FORMAT_ERROR
	}
	return idx + len(SEPERATOR), r, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := newRequest()
	buff := make([]byte, 1024)
	nBuff:= 0

	for !request.isDone(){
		n, err := reader.Read(buff[nBuff:])
		if err != nil {
			return nil, err
		}
		nBuff += n
		n, err = request.parse(buff[:nBuff])
		// EOF is also an error we don't wanna break it since it expected
		if err != nil && err != io.EOF {
			return nil, err
		}
		copy(buff, buff[n:nBuff])
		nBuff -= n
	}
	return request, nil
}