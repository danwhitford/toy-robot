package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/danwhitford/toyrobot/toyrobot"
)

func main() {
	r := toyrobot.NewRobot()
	buf := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			continue
		}
		err = r.ReadInstruction(string(line))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
