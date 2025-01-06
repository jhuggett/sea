package jsonrpc

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type JSONRPC interface {
	Send(method string, params interface{}) (chan Response, error)
	Receive(method string, callback func(params json.RawMessage) (interface{}, error)) func()
}

type jsonrpc struct {
	conn *websocket.Conn

	requestRegistryMutex sync.Mutex
	requestRegistry      map[string]func(params json.RawMessage) (interface{}, error)

	responseRegistryMutex sync.Mutex
	responseRegistry      map[string]chan Response

	ClosedChan chan struct{}
}

type connMessage struct {
	ID string `json:"id"`

	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`

	Result json.RawMessage `json:"result,omitempty"`
	Error  json.RawMessage `json:"error,omitempty"`
}

func New(conn *websocket.Conn) *jsonrpc {
	rpc := &jsonrpc{
		conn:             conn,
		requestRegistry:  make(map[string]func(params json.RawMessage) (interface{}, error)),
		responseRegistry: make(map[string]chan Response, 1),
	}

	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var msg connMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				slog.Warn("jsonrpc: failed to unmarshal message: %v", err)
				continue
			}

			if msg.Method != "" {
				rpc.requestRegistryMutex.Lock()
				if callback, ok := rpc.requestRegistry[msg.Method]; ok {
					rpc.requestRegistryMutex.Unlock()
					result, err := callback(msg.Params)

					var respMessage []byte

					if err != nil {
						slog.Error("jsonrpc: error in callback: %v", err)

						marshaledErr, err := json.Marshal(err)
						if err != nil {
							slog.Warn("jsonrpc: failed to marshal error: %v", err)
							continue
						}

						data, err := json.Marshal(Response{
							ID:    msg.ID,
							Error: marshaledErr,
						})
						if err != nil {
							slog.Warn("jsonrpc: failed to marshal response: %v", err)
							continue
						}
						respMessage = data
					} else {
						marshaledResult, err := json.Marshal(result)
						if err != nil {
							slog.Warn("jsonrpc: failed to marshal result", "err", err, "result", result)
							continue
						}

						data, err := json.Marshal(Response{
							ID:     msg.ID,
							Result: marshaledResult,
						})
						if err != nil {
							slog.Warn("jsonrpc: failed to marshal response: %v", err)
							continue
						}
						respMessage = data
					}

					err = conn.WriteMessage(websocket.TextMessage, respMessage)
					if err != nil {
						slog.Warn("jsonrpc: failed to write message: %v", err)
					}
				}
			} else {
				rpc.responseRegistryMutex.Lock()
				if respChannel, ok := rpc.responseRegistry[msg.ID]; ok {
					respChannel <- Response{
						ID:     msg.ID,
						Result: msg.Result,
						Error:  msg.Error,
					}
					delete(rpc.responseRegistry, msg.ID)
				}
				rpc.responseRegistryMutex.Unlock()
			}
		}

		// rpc.ClosedChan <- struct{}{}
	}()

	return rpc
}

type Request struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type Response struct {
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  json.RawMessage `json:"error,omitempty"`
}

func (j *jsonrpc) Send(method string, params interface{}) (chan Response, error) {
	req := Request{
		ID:     uuid.New().String(),
		Method: method,
		Params: params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	j.conn.WriteMessage(websocket.TextMessage, data)

	respChannel := make(chan Response, 1)

	j.responseRegistryMutex.Lock()
	j.responseRegistry[req.ID] = respChannel
	j.responseRegistryMutex.Unlock()

	return respChannel, nil
}

func (j *jsonrpc) Receive(method string, callback func(params json.RawMessage) (interface{}, error)) func() {
	j.requestRegistry[method] = callback

	return func() {
		delete(j.requestRegistry, method)
	}
}
