package model

import (
	"errors"
	"testing"

	"github.com/c2fo/testify/require"
)

func Test_ConnError(t *testing.T) {
	req := require.New(t)

	tCases := []struct {
		connId      string
		action      string
		code        int
		text        string
		errs        []interface{}
		requiredErr string
	}{
		{
			connId: "conn-1",
			action: "connect",
			code:   -1,
			errs: []interface{}{
				"Failed to connect",
			},
			text:        "Failed to connect",
			requiredErr: "ConnID: conn-1 Action: connect ErrorCode: -1 Error: Failed to connect",
		},
		{
			connId: "conn-2",
			action: "read raw",
			code:   -2,
			errs: []interface{}{
				errors.New("Failed to read raw data from WebSocket"),
			},
			text:        "Failed to read raw data from WebSocket",
			requiredErr: "ConnID: conn-2 Action: read raw ErrorCode: -2 Error: Failed to read raw data from WebSocket",
		},
		{
			connId: "conn-3",
			action: "send prepared message",
			code:   -3,
			errs: []interface{}{
				"Failed to send prepared message on the WebSocket",
				errors.New("Write error 420"),
			},
			text:        "Failed to send prepared message on the WebSocket: Write error 420",
			requiredErr: "ConnID: conn-3 Action: send prepared message ErrorCode: -3 Error: Failed to send prepared message on the WebSocket: Write error 420",
		},
		{
			connId:      "conn-4",
			action:      "",
			code:        0,
			errs:        []interface{}{},
			text:        "",
			requiredErr: "ConnID: conn-4 Action:  ErrorCode: 0 Error: ",
		},
	}

	for _, tCase := range tCases {
		err := NewConnError(tCase.connId, tCase.action, tCase.code, tCase.errs...)
		req.Error(err, "NewConnError should create an Error")
		req.IsType(err, &ConnError{}, "NewConnError should create a ConnError")
		req.Equal(tCase.connId, err.ConnId(), "err.ConnId should return the Connection ID")
		req.Equal(tCase.action, err.Action(), "err.Action should return the intended action")
		req.Equal(tCase.code, err.Code(), "err.Code should return the error code")
		req.Equal(tCase.text, err.Text(), "err.Text should return the error text")
		req.Equal(tCase.requiredErr, err.Error(), "The error must be in the 'Conn ID: %s Action: %s Error Code: %d Error: %s' format")
	}
}
