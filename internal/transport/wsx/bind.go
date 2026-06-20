// Package wsx
package wsx

import (
	"encoding/json"

	"github.com/AtoriUzawa/cira"
)

// BindJSON automatically sends a response on bind error
func BindJSON[T any](c *cira.Context, t *T) bool {
	err := json.Unmarshal(c.Message.Data, t)
	if err != nil {
		c.Resp(map[string]string{"msg": "invalid"})
		return false
	}

	return true
}
