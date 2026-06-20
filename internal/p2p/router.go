package p2p

import (
	"github.com/AtoriUzawa/cira"
	"github.com/gin-gonic/gin"
)

// RegisterRouter registers HTTP routes for P2P operations.
func RegisterRouter(r *gin.RouterGroup, h *Handler) {
	api := r.Group("/p2p")
	{
		api.POST("/list", h.List)
	}
}

// RegisterWSRouter registers WebSocket routes for P2P operations.
func RegisterWSRouter(r *cira.RouterGroup, h *WSHandler) {
	r = r.Group("p2p")
	{
		r.On("register", h.Register)
		r.On("unregister", h.Unregister)
		r.On("call", h.Call)
		r.On("hangup", h.Hangup)
	}
}
