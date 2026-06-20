package live

import (
	"github.com/AtoriUzawa/cira"
	"github.com/gin-gonic/gin"
)

// RegisterRouter registers HTTP routes for live streaming operations.
func RegisterRouter(r *gin.RouterGroup, h *Handler) {
	api := r.Group("/live")
	{
		api.POST("/list", h.List)
	}
}

// RegisterWSRouter registers WebSocket routes for live streaming operations.
func RegisterWSRouter(r *cira.RouterGroup, h *WSHandler) {
	r = r.Group("live")
	{
		r.On("create", h.Create)
		r.On("join", h.Join)
		r.On("leave", h.Leave)
	}
}
