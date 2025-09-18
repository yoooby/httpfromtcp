package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
)


type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError


func WriteHandlerError(w io.Writer, err HandlerError) error {
    // First, write the status line
    if writeErr := response.WriteStatusLine(w, err.StatusCode); writeErr != nil {
        return writeErr
    }

	// write headers
	headers := response.GetDefaultHeaders(len(err.Message))
	writeErr := response.WriteHeaders(w,headers)
	if writeErr != nil {
		return writeErr
	}

    // Then write the error message in the body
    if err.Message != "" {
        _, writeErr = fmt.Fprintf(w, "%s\r\n", err.Message)
        return writeErr
    }

    return nil
}

