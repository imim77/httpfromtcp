package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
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
				str += string(data[:i]) // up to index i
				data = data[i+1:]       // + 1 because we don't want that new line
				ch <- str
				str = ""
			}

			str += string(data)
		}

		if len(str) != 0 {
			ch <- str
		}

	}() // anonymus function, it executes immediatly
	return ch

}

func main() {
	f, err := os.Open("message.txt")
	if err != nil {
		log.Fatal("errror", "error", err)
	}
	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
