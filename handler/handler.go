package handler

import "swiftdb/resp"

// commands handler

// writing PING command
func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

	return resp.Value{Typ: "string", Str: args[0].Bulk}
}

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
}
