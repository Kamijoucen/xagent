package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		handleLine(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		write(response{JSONRPC: "2.0", ID: nil, Error: &rpcError{Code: -32603, Message: err.Error()}})
	}
}

func handleLine(line []byte) {
	var req request
	if err := json.Unmarshal(line, &req); err != nil {
		write(response{JSONRPC: "2.0", ID: nil, Error: &rpcError{Code: -32700, Message: "invalid JSON"}})
		return
	}
	if req.Method != "execute" {
		write(response{JSONRPC: "2.0", ID: req.ID, Error: &rpcError{Code: -32601, Message: "method not found"}})
		return
	}
	write(response{JSONRPC: "2.0", ID: req.ID, Result: map[string]interface{}{"content": string(req.Params), "success": true}})
}

func write(resp response) {
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintln(os.Stdout, `{"jsonrpc":"2.0","id":null,"error":{"code":-32603,"message":"marshal error"}}`)
		return
	}
	fmt.Fprintln(os.Stdout, string(data))
}
