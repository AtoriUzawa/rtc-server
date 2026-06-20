package p2p

import (
	"encoding/json"

	"github.com/AtoriUzawa/vlink-server/internal/transport/wsx"
	"github.com/AtoriUzawa/cira"
)

// WSHandler handles WebSocket events for the P2P module.
type WSHandler struct {
	s *Service
}

// NewWSHandler creates a new WSHandler with the given Service.
func NewWSHandler(s *Service) *WSHandler {
	return &WSHandler{
		s: s,
	}
}

// Register handles the peer registration WebSocket event.
func (h *WSHandler) Register(c *cira.Context) {
	var req RegisterReq
	err := json.Unmarshal(c.Message.Data, &req)
	if err != nil {
		return
	}

	c.Conn.OnClose(func() {
		h.s.m.Unregister(req.ID)
	})

	h.s.m.Register(&Peer{
		ID: req.ID,
	})
	wsx.OK(c)
}

// Unregister handles the peer unregistration WebSocket event.
func (h *WSHandler) Unregister(c *cira.Context) {
	var req UnregisterReq
	err := json.Unmarshal(c.Message.Data, &req)
	if err != nil {
		return
	}

	h.s.m.Unregister(req.ID)
	wsx.OK(c)
}

// Call handles the call initiation WebSocket event.
func (h *WSHandler) Call(c *cira.Context) {
	var req CallReq
	err := json.Unmarshal(c.Message.Data, &req)
	if err != nil {
		return
	}

	h.s.Call(c, &req)
}

// Hangup handles the hangup WebSocket event.
func (h *WSHandler) Hangup(c *cira.Context) {
	var req HangupReq
	err := json.Unmarshal(c.Message.Data, &req)
	if err != nil {
		return
	}

	h.s.Hangup(&req)
}
