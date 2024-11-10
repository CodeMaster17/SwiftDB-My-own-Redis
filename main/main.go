package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"swiftdb/resp"
)

func main() {
	// TCP listener so that client can communicate with it
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	// closing connection afer finishing
	defer conn.Close()

	// creating infinite loop to receive commands from client and respond to them
	for {
		resp := resp.NewResp(conn)

		// read message from client
		value, err := resp.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from the client: ", err.Error())
			os.Exit(1)
		}

		fmt.Println(value)

		// ignore request and send back pong
		conn.Write([]byte("+OK\r\n"))
	}
}
