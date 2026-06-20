package meeting

import (
	"github.com/AtoriUzawa/vlink-server/internal/transport/wsx"
	"github.com/AtoriUzawa/cira"
)

// WSHandler handles WebSocket events for the meeting module.
type WSHandler struct {
	s *Service
}

// NewWSHandler creates a new WSHandler with the given Service.
func NewWSHandler(s *Service) *WSHandler {
	return &WSHandler{s: s}
}

// Create handles the create meeting room WebSocket event.
func (h *WSHandler) Create(c *cira.Context) {
	var req CreateReq
	if !wsx.BindJSON(c, &req) {
		return
	}

	err := h.s.Create(&req, c.Conn)
	if err != nil {
		wsx.FailedWithErr(c, err)
		return
	}

	wsx.OK(c)
}

// Join handles the join meeting room WebSocket event.
func (h *WSHandler) Join(c *cira.Context) {
	var req JoinReq
	if !wsx.BindJSON(c, &req) {
		return
	}

	err := h.s.Join(&req, c.Conn)
	if err != nil {
		wsx.FailedWithErr(c, err)
		return
	}

	wsx.OK(c)
}

// Leave handles the leave meeting room WebSocket event.
func (h *WSHandler) Leave(c *cira.Context) {
	var req LeaveReq
	if !wsx.BindJSON(c, &req) {
		return
	}

	err := h.s.Leave(&req)
	if err != nil {
		wsx.FailedWithErr(c, err)
		return
	}

	wsx.OK(c)
}
