package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	curLine := ""
	for {
		data := make([]byte, 8)
		_, err := file.Read(data)
		if err != nil {
			break
		}
		data_s := string(data)
		strings := strings.Split(data_s, "\n")
		curLine += strings[0]
		if len(strings) > 1 {
			fmt.Printf("read: %s\n", curLine)
			curLine = strings[1]
		}
	}

}
