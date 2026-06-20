// Package httpx
package httpx

import (
	"errors"
	"net/http"

	"github.com/AtoriUzawa/vlink-server/pkg/xerror"
	"github.com/gin-gonic/gin"
)

// Response represents a standard JSON API response with code, message and data fields.
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

const (
	// MsgSuccess is the default success message returned in API responses.
	MsgSuccess = "success"
	// MsgInternalErr is the default internal error message returned in API responses.
	MsgInternalErr = "internal error"
)

// OkWithData sends a 200 JSON response with the given data payload.
func OkWithData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, &Response{
		Code: http.StatusOK,
		Msg:  MsgSuccess,
		Data: data,
	})
}

// Ok sends a 200 JSON success response with an empty data object.
func Ok(c *gin.Context) {
	c.JSON(http.StatusOK, &Response{
		Code: http.StatusOK,
		Msg:  MsgSuccess,
		Data: struct{}{},
	})
}

// Err sends a JSON error response derived from the given error.
// If the error is an *xerror.XError, its code and message are used;
// otherwise a 500 internal server error is returned.
func Err(c *gin.Context, err error) {
	var xe *xerror.XError
	if ok := errors.As(err, &xe); !ok {
		c.JSON(http.StatusInternalServerError, &Response{
			Code: http.StatusInternalServerError,
			Msg:  MsgInternalErr,
			Data: nil,
		})
		return
	}

	c.JSON(xe.Code, &Response{
		Code: xe.Code,
		Msg:  xe.Msg,
		Data: nil,
	})
}
