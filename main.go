package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/danwhitford/toyrobot/toyrobot"
)

func main() {
	r := toyrobot.NewRobot()
	buf := bufio.NewReader(os.Stdin)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			continue
		}
		err = r.RunProgram(string(line))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
