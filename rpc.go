package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const JSONRPC_VERSION = "2.0"

type Request struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	Id      any             `json:"id,omitempty"`
}

type Response struct {
	Version string `json:"jsonrpc"`
	Result  any    `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
	Id      any    `json:"id,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Message)
}

func UnmarshalRequest(data []byte) (*Request, error) {
	var request Request
	err := json.Unmarshal(data, &request)
	if err != nil {
		return nil, &Error{
			Code:    -32700,
			Message: err.Error(),
		}
	}
	if request.Version != JSONRPC_VERSION {
		return nil, &Error{
			Code:    -32600,
			Message: "missing or invalid jsonrpc version",
		}
	}
	if request.Method == "" {
		return nil, &Error{
			Code:    -32600,
			Message: "missing method",
		}
	}
	return &request, nil
}

func UnmarshalParams[T any](data json.RawMessage) (T, error) {
	var result T
	err := json.Unmarshal([]byte(data), &result)
	return result, err
}

func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func NewResponse(id any, result any, err error) *Response {
	if v, ok := result.(*Response); ok {
		return v
	}
	if err == nil {
		return &Response{
			Version: JSONRPC_VERSION,
			Result:  result,
			Error:   nil,
			Id:      id,
		}
	}
	if v, ok := err.(*Error); ok {
		return &Response{
			Version: JSONRPC_VERSION,
			Result:  nil,
			Error:   v,
			Id:      id,
		}
	}
	return &Response{
		Version: JSONRPC_VERSION,
		Result:  nil,
		Error: &Error{
			Code:    -32603,
			Message: err.Error(),
		},
		Id: id,
	}
}

type Method func(id any, params json.RawMessage) (any, error)

func HTTPError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}

type Handler struct {
	Methods map[string]Method
}

func (h *Handler) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HTTPError(w, http.StatusBadRequest, fmt.Errorf("Wrong HTTP method"))
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, err)
		return
	}
	req, err := UnmarshalRequest(body)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, err)
		return
	}
	m, ok := h.Methods[req.Method]
	if !ok {
		HTTPError(w, http.StatusBadRequest, fmt.Errorf("Unknown jsonrpc method"))
		return
	}
	result, err := m(req.Id, req.Params)
	resp := NewResponse(req.Id, result, err)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, err)
		return
	}
	resp_data, err := resp.Marshal()
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(resp_data)
}
