package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"net"
)


type Server struct {
	listener net.Listener
	handler Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{listener: listener, handler: handler}
	go s.listen() // start listening in the background
	return s, nil
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Error Accepting connection: %v\n", err)
			continue
		}
		go s.handle(conn) // handle each connection in its own goroutine
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
/* 
	headers := response.GetDefaultHeaders(len(body))
	response.WriteStatusLine(conn, response.StatusCode200)
	response.WriteHeaders(conn, headers)
	conn.Write([]byte(body)) */

	r, parseError := request.RequestFromReader(conn)
	if parseError != nil {
		WriteHandlerError(conn, HandlerError{StatusCode: 500, Message: parseError.Error()})
	}
	var buf bytes.Buffer

	err := s.handler(&buf, r)
	if err != nil {
		WriteHandlerError(conn, *err)
		return
	}
	// no erros
	body := buf.Bytes()
	headers := response.GetDefaultHeaders(len(body))

	response.WriteStatusLine(conn, response.StatusCode200)
	response.WriteHeaders(conn, headers)
	conn.Write(body)
}