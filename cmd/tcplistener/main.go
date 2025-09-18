package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func mainHandler(w io.Writer, r *request.Request) *server.HandlerError {
	if r.RequestLine.RequestTarget == "/bad" {
		return &server.HandlerError{
			StatusCode: 400,
			Message: "Waaaaaaaa3" ,
		}
	} else if r.RequestLine.RequestTarget == "/fuck" {
		return &server.HandlerError{
			StatusCode: 500,
			Message: "PANIC" ,
		}  
	} else {
		w.Write([]byte("ALL GOOD YEAAAAH"))
	}
	return nil
}

func main() {
	server, err := server.Serve(port, mainHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}