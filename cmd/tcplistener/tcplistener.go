package main

/* func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("Error binding port")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error accepting")
		}
		log.Printf("Accepted Connection\n")
		r, err := request.RequestFromReader(conn)
		if err != nil {
			conn.Close()
			log.Fatal("error")
		}
		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for k,v := range r.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Printf("Body:\n %s\n", r.Body)
		resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 2\r\n\r\nHello World!\r\n"
		conn.Close()
		log.Printf("connection closed\n")
	}
} */