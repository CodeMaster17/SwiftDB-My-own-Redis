package main

import (
	"fmt"
	"net"
	"strings"
	"swiftdb/aof"
	"swiftdb/handler"
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

	// aof
	aof, err := aof.NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(value resp.Value) {
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

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

		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		/*
			Value{
				typ: "array",
				array: []Value{
					Value{typ: "bulk", bulk: "SET"},
					Value{typ: "bulk", bulk: "name"},
					Value{typ: "bulk", bulk: "Ahmed"},
				},
			}
		*/

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		/*
			command := Value{typ: "bulk", bulk: "SET"}.bulk // "SET"

			args := []Value{
				Value{typ: "bulk", bulk: "name"},
				Value{typ: "bulk", bulk: "Ahmed"},
			}
		*/

		writer := resp.NewWriter(conn)

		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(resp.Value{Typ: "string"})
			continue
		}

		// if command == "SET" || command == "HSET" {
		// }
		aof.Write(value)

		// with this to add error handling and logging
		if err := aof.Write(value); err != nil {
			fmt.Println("Failed to write to AOF:", err)
		} else {
			fmt.Println("AOF Write successful for command:", command)
		}

		result := handler(args)
		writer.Write(result)
	}
}
