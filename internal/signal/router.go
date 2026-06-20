package signal

import (
	"github.com/AtoriUzawa/cira"
)

// RegisterWSRouter registers WebSocket routes for signal operations.
func RegisterWSRouter(r *cira.RouterGroup, h *WSHandler) {
	r = r.Group("signal.rtc")
	{
		r.On("register", h.Register)
		r.On("unregister", h.UnRegister)
		r.On("forward", h.Signal)
	}
}
