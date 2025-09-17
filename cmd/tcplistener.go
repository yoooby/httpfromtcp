package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func(){
		defer f.Close()
		defer close(out)

		str := ""
		for {
			data := make([]byte, 8)
			_, err := f.Read(data)
			if err != nil {
				break
			}

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				out <- str
				str = ""
			} 
			str += string(data)
		}

		
		if len(str) != 0 {
			out <- str
		}
	}()



	return out
}


func main() {

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
		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
		conn.Close()
		log.Printf("connection closed\n")


	}
}