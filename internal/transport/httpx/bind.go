package httpx

import (
	"github.com/AtoriUzawa/vlink-server/pkg/xerror"
	"github.com/gin-gonic/gin"
)

// BindJSON automatically sends a response on bind error
func BindJSON[T any](c *gin.Context, t *T) bool {
	err := c.ShouldBindJSON(t)
	if err != nil {
		Err(c, xerror.WithMsg("invalid param", xerror.Wrap(err, "failed to bind json")))
		return false
	}

	return true
}
