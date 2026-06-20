package meeting

import (
	"github.com/AtoriUzawa/cira"
	"github.com/gin-gonic/gin"
)

// RegisterRouter registers HTTP routes for meeting operations.
func RegisterRouter(r *gin.RouterGroup, h *Handler) {
	api := r.Group("/meeting")
	{
		api.POST("/list", h.List)
	}
}

// RegisterWSRouter registers WebSocket routes for meeting operations.
func RegisterWSRouter(r *cira.RouterGroup, h *WSHandler) {
	r = r.Group("meeting")
	{
		r.On("create", h.Create)
		r.On("join", h.Join)
		r.On("leave", h.Leave)
	}
}
