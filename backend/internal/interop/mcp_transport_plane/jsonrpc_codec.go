package mcptransportplane

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const maxJSONRPCBytes = 4 << 20

// JSONRPCRequest is a strict JSON-RPC 2.0 request.
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCNotification is a strict JSON-RPC 2.0 notification with no ID.
type JSONRPCNotification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCError is a strict JSON-RPC error object.
type JSONRPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// JSONRPCResponse is a strict JSON-RPC 2.0 response.
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// EncodeRequest validates and encodes one request.
func EncodeRequest(request JSONRPCRequest) ([]byte, error) {
	if err := validateRequest(request); err != nil {
		return nil, err
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("encode json-rpc request: %w", err)
	}
	if len(data) > maxJSONRPCBytes {
		return nil, errors.New("json-rpc request exceeds size limit")
	}
	return data, nil
}

// EncodeNotification validates and encodes one JSON-RPC notification.
func EncodeNotification(notification JSONRPCNotification) ([]byte, error) {
	if err := validateNotification(notification); err != nil {
		return nil, err
	}
	data, err := json.Marshal(notification)
	if err != nil {
		return nil, fmt.Errorf("encode json-rpc notification: %w", err)
	}
	if len(data) > maxJSONRPCBytes {
		return nil, errors.New("json-rpc notification exceeds size limit")
	}
	return data, nil
}

// DecodeRequest strictly decodes one request.
func DecodeRequest(data []byte) (JSONRPCRequest, error) {
	if len(data) == 0 || len(data) > maxJSONRPCBytes {
		return JSONRPCRequest{}, errors.New("json-rpc request size is invalid")
	}
	var request JSONRPCRequest
	if err := decodeStrictJSON(data, &request); err != nil {
		return JSONRPCRequest{}, fmt.Errorf("decode json-rpc request: %w", err)
	}
	if err := validateRequest(request); err != nil {
		return JSONRPCRequest{}, err
	}
	return request, nil
}

// DecodeResponse strictly decodes one response and enforces result/error
// exclusivity.
func DecodeResponse(data []byte) (JSONRPCResponse, error) {
	if len(data) == 0 || len(data) > maxJSONRPCBytes {
		return JSONRPCResponse{}, errors.New("json-rpc response size is invalid")
	}
	var response JSONRPCResponse
	if err := decodeStrictJSON(data, &response); err != nil {
		return JSONRPCResponse{}, fmt.Errorf("decode json-rpc response: %w", err)
	}
	if response.JSONRPC != "2.0" || !validJSONRPCID(response.ID) {
		return JSONRPCResponse{}, errors.New("json-rpc response envelope is invalid")
	}
	hasResult := len(response.Result) > 0
	hasError := response.Error != nil
	if hasResult == hasError {
		return JSONRPCResponse{}, errors.New("json-rpc response requires exactly one result or error")
	}
	if hasError && (response.Error.Code == 0 || response.Error.Message == "") {
		return JSONRPCResponse{}, errors.New("json-rpc error object is invalid")
	}
	return response, nil
}

func validateRequest(request JSONRPCRequest) error {
	if request.JSONRPC != "2.0" ||
		!validJSONRPCID(request.ID) ||
		!validJSONRPCMethod(request.Method) {
		return errors.New("json-rpc request envelope is invalid")
	}
	if len(request.Params) > 0 && !json.Valid(request.Params) {
		return errors.New("json-rpc params are invalid")
	}
	return nil
}

func validateNotification(notification JSONRPCNotification) error {
	if notification.JSONRPC != "2.0" ||
		!validJSONRPCMethod(notification.Method) {
		return errors.New("json-rpc notification envelope is invalid")
	}
	if len(notification.Params) > 0 && !json.Valid(notification.Params) {
		return errors.New("json-rpc params are invalid")
	}
	return nil
}

func validJSONRPCMethod(method string) bool {
	return method != "" && len(method) <= 128
}

func validJSONRPCID(id json.RawMessage) bool {
	if len(id) == 0 || bytes.Equal(id, []byte("null")) {
		return false
	}
	var stringID string
	if json.Unmarshal(id, &stringID) == nil {
		return stringID != "" && len(stringID) <= 128
	}
	var numberID int64
	return json.Unmarshal(id, &numberID) == nil
}

func decodeStrictJSON(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("json contains trailing data")
	}
	return nil
}
