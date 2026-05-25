package plugins

import "encoding/json"

// Request 是 JSON-RPC 2.0 请求结构。
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response 是 JSON-RPC 2.0 响应结构。
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error 是 JSON-RPC 2.0 错误结构。
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewResultResponse 创建成功响应。
func NewResultResponse(id interface{}, result interface{}) Response {
	return Response{JSONRPC: "2.0", ID: id, Result: result}
}

// NewErrorResponse 创建错误响应。
func NewErrorResponse(id interface{}, code int, message string) Response {
	return Response{JSONRPC: "2.0", ID: id, Error: &Error{Code: code, Message: message}}
}
