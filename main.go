package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close()
		defer close(ch)

		curLine := ""
		for {
			data := make([]byte, 8)
			_, err := f.Read(data)
			if err != nil {
				break
			}
			data_s := string(data)
			strings := strings.Split(data_s, "\n")
			curLine += strings[0]
			if len(strings) > 1 {
				ch <- curLine
				curLine = strings[1]
			}
		}
	}()
	return ch

}
func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
