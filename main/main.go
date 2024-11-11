package main

import (
	"fmt"
	"net"
	"swiftdb/resp"
)

func main() {
	fmt.Println("Listening on port :6379")

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
		response := resp.NewResp(conn)
		value, err := response.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		_ = value

		writer := resp.NewWriter(conn)
		writer.Write(resp.NewValue("string", "OK", 0, "", nil))
	}
}
