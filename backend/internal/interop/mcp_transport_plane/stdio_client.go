package mcptransportplane

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

// StdioClient exchanges newline-delimited JSON-RPC over a pipe created by
// Layer 6. It never spawns a process.
type StdioClient struct {
	mu     sync.Mutex
	pipe   io.ReadWriteCloser
	reader *bufio.Reader
}

// NewStdioClient wraps an already-sandboxed pipe.
func NewStdioClient(pipe io.ReadWriteCloser) (*StdioClient, error) {
	if pipe == nil {
		return nil, errors.New("mcp stdio pipe is required")
	}
	return &StdioClient{
		pipe:   pipe,
		reader: bufio.NewReaderSize(pipe, 64<<10),
	}, nil
}

// Call performs one serialized request/response exchange.
func (c *StdioClient) Call(
	ctx context.Context,
	request JSONRPCRequest,
) (JSONRPCResponse, error) {
	if c == nil || c.pipe == nil || c.reader == nil {
		return JSONRPCResponse{}, errors.New("mcp stdio client is not configured")
	}
	if err := ctx.Err(); err != nil {
		return JSONRPCResponse{}, fmt.Errorf("mcp stdio call: %w", err)
	}
	data, err := EncodeRequest(request)
	if err != nil {
		return JSONRPCResponse{}, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if _, err := c.pipe.Write(append(data, '\n')); err != nil {
		return JSONRPCResponse{}, fmt.Errorf("write mcp stdio request: %w", err)
	}
	response, err := c.reader.ReadBytes('\n')
	if err != nil {
		return JSONRPCResponse{}, fmt.Errorf("read mcp stdio response: %w", err)
	}
	if len(response) > maxJSONRPCBytes+1 {
		return JSONRPCResponse{}, errors.New("mcp stdio response exceeds size limit")
	}
	return DecodeResponse(response)
}

// Close closes only the sanctioned pipe. Layer 6 remains responsible for the
// sandboxed process lifecycle.
func (c *StdioClient) Close() error {
	if c == nil || c.pipe == nil {
		return nil
	}
	return c.pipe.Close()
}
