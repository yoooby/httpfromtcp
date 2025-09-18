package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)


type StatusCode int

const (
	StatusCode200 StatusCode = 200;
	StatusCode400 StatusCode = 400;
	StatusCode500 StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var reason string
	switch statusCode {
	case StatusCode200:
		reason = "OK"
	case StatusCode400:
		reason = "Bad Request"
	case StatusCode500:	
		reason = "Internal Server Error"
	default:
		reason = ""
	}
	if reason != "" {
		_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", statusCode, reason)
		return err
	}
	_, err := fmt.Fprintf(w, "HTTP/1.1 %d\r\n", statusCode)
	return err
}

func GetDefaultHeaders(contentLength int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = fmt.Sprint(contentLength)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k,v := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", k, v)
		if err != nil {
			return fmt.Errorf("Error occured when writing %s : %s | %v", k, v, err)
		}
	}
	_, err := fmt.Fprintf(w, "\r\n")
	return err
}