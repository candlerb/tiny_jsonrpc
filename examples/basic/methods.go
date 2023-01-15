package main

import (
	"encoding/json"
	"github.com/candlerb/tiny_jsonrpc"
)

var MyHandler = &rpc.Handler{
	Methods: map[string]rpc.Method{
		"ping": pong,
		"add":  add,
	},
}

func pong(id any, params json.RawMessage) (any, error) {
	return "pong", nil
}

// Each method is required to parse its own params from a json.RawMessage
func add(id any, params json.RawMessage) (any, error) {
	args, err := rpc.UnmarshalParams[[]float64](params)
	if err != nil {
		return nil, &rpc.Error{
			Code:    -32602,
			Message: err.Error(),
		}
	}
	if len(args) != 2 {
		return nil, &rpc.Error{
			Code:    -32602,
			Message: "Invalid number of params",
		}
	}
	res := args[0] + args[1]
	return res, nil
}
