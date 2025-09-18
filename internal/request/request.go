package request

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"log"
	"slices"
	"strconv"
)


type RequestState int

const (
	RequestInitialized RequestState = iota
	RequestParsingHeaders
	RequestParsingBody
	RequestDone
)

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	Body []byte
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
			r.State = RequestParsingHeaders
		case RequestParsingHeaders:
			if r.Headers == nil {
				r.Headers = headers.NewHeaders()
			}
			n, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			read += n
			if done {
				if r.Headers.GET("Content-Length") != "" {
					r.State = RequestParsingBody
				} else {
					r.State = RequestDone
				}
			}
		case RequestParsingBody:
			n, done, err := r.parseBody(data[read:])
			if err != nil {
				return read, err // propagate bytes consumed so far
			}
			read += n
			if done {
				r.State = RequestDone
				break outer
			}
			if n == 0 {
				break outer // need more data, keep the state as-is
			}


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
func (r *Request) parseBody(data []byte) (int, bool, error) {
	cl := r.Headers.GET("Content-Length")
	// it shouldnt be an issue since we making sure it fucking exists
	contentLength, err := strconv.Atoi(cl)
	if err != nil {
		return 0, false, fmt.Errorf("invalid Content-Length: %v", err)
	}

	remaining := contentLength - len(r.Body)
	if remaining <= 0 {
		// Body already complete
		if len(data) > 0 {
			return 0, true, fmt.Errorf("body exceeds Content-Length")
		}
		return 0, true, nil
	}

	toRead := min(len(data), remaining)
	r.Body = append(r.Body, data[:toRead]...)

	// If there's still extra data after finishing the body, that's an error
	if len(r.Body) == contentLength && len(data) > toRead {
		return toRead, true, fmt.Errorf("body exceeds Content-Length")
	}

	done := len(r.Body) == contentLength
	return toRead, done, nil
}




func RequestFromReader(reader io.Reader) (*Request, error) {

	request := newRequest()
	buff := make([]byte, 1024)
	nBuff:= 0

	for !request.isDone(){
		n, err := reader.Read(buff[nBuff:])
		if err != nil{
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