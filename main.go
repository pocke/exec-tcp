package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			cmd := make([]string, 0)
			err := json.NewDecoder(conn).Decode(&cmd)
			if err != nil {
				log.Println(err)
				return
			}

			fmt.Println("> " + strings.Join(cmd, " "))
			c := exec.Command(cmd[0], cmd[1:]...)
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
			c.Stdin = os.Stdin
			if err := c.Run(); err != nil {
				log.Println(err)
				return
			}
		}(conn)
	}
}
