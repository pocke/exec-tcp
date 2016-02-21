package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"reflect"
)

func getBodyFromVimChan(ch interface{}) (string, error) {
	s := reflect.ValueOf(ch)
	if s.Kind() != reflect.Slice {
		return "", fmt.Errorf("%v is not slice", ch)
	}
	return s.Index(1).Interface().(string), nil
}

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("> ")
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		if err := Exec(conn); err != nil {
			log.Println(err)
		}
		fmt.Print("> ")
	}
}

func Exec(r io.ReadCloser) error {
	defer r.Close()

	v := make([]interface{}, 2)
	err := json.NewDecoder(r).Decode(&v)
	if err != nil {
		return err
	}

	cmd, err := getBodyFromVimChan(v)
	if err != nil {
		return err
	}

	fmt.Println(cmd)
	c := exec.Command("bash", "-c", cmd)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}
