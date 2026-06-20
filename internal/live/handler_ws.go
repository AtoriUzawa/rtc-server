package live

import (
	"github.com/AtoriUzawa/vlink-server/internal/transport/wsx"
	"github.com/AtoriUzawa/cira"
)

// WSHandler handles WebSocket events for the live module.
type WSHandler struct {
	s *Service
}

// NewWSHandler creates a new WSHandler with the given Service.
func NewWSHandler(s *Service) *WSHandler {
	return &WSHandler{s: s}
}

// Create handles the create live room WebSocket event.
func (h *WSHandler) Create(c *cira.Context) {
	var req CreateReq
	if !wsx.BindJSON(c, &req) {
		return
	}

	if err := h.s.Create(&req, c.Conn); err != nil {
		wsx.FailedWithErr(c, err)
		return
	}

	wsx.OK(c)
}

// Join handles the join live room WebSocket event.
func (h *WSHandler) Join(c *cira.Context) {
	var req JoinReq
	if !wsx.BindJSON(c, &req) {
		return
	}

	if err := h.s.Join(&req, c.Conn); err != nil {
		wsx.FailedWithErr(c, err)
		return
	}

	wsx.OK(c)
}

// Leave handles the leave live room WebSocket event.
func (h *WSHandler) Leave(c *cira.Context) {
	var req LeaveReq
	if !wsx.BindJSON(c, &req) {
		return
	}

	if err := h.s.Leave(&req); err != nil {
		wsx.FailedWithErr(c, err)
		return
	}

	wsx.OK(c)
}
