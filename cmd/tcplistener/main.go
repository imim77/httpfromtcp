package main

import (
	"fmt"
	"log"
	"net"

	"github.com/imim77/httpfromtcp/internal/request"
)

const port = ":42069"

/*func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(ch)
		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}
			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				ch <- str
				str = ""
			}
			str += string(data)
		}
		if len(str) != 0 {
			ch <- str
		}
	}()
	return ch

}
*/

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP: %s\n", err.Error())
	}
	defer listener.Close()
	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Connection has been accepted")

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
	}
}
