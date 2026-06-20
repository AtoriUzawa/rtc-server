// Package wsx provides WebSocket response helpers.
package wsx

import "github.com/AtoriUzawa/cira"

// OK sends a successful WebSocket response with "msg": "ok".
func OK(c *cira.Context) {
	c.Resp(map[string]string{"msg": "ok"})
}

// Failed sends a failure WebSocket response with "msg": "failed".
func Failed(c *cira.Context) {
	c.Resp(map[string]string{"msg": "failed"})
}

// FailedWithErr sends a failure WebSocket response with the error message.
func FailedWithErr(c *cira.Context, err error) {
	c.Resp(map[string]string{"msg": err.Error()})
}
