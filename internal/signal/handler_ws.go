package signal

import (
	"encoding/json"
	"log"

	"github.com/AtoriUzawa/vlink-server/internal/transport/wsx"
	"github.com/AtoriUzawa/cira"
)

// WSHandler handles WebSocket events for signal forwarding.
type WSHandler struct {
	m *Manager
}

// NewWSHandler creates a new WSHandler with the given Manager.
func NewWSHandler(m *Manager) *WSHandler {
	return &WSHandler{m: m}
}

// Signal forwards a WebRTC signaling message to the target peer.
func (h *WSHandler) Signal(c *cira.Context) {
	log.Printf(
		"[SIGNAL] recv route=%s",
		c.Message.Route,
	)

	var req SignalMessage

	err := json.Unmarshal(c.Message.Data, &req)
	if err != nil {
		log.Printf(
			"[SIGNAL] unmarshal failed: %v\nraw=%s",
			err,
			string(c.Message.Data),
		)

		wsx.FailedWithErr(c, err)

		return
	}

	log.Printf(
		"[SIGNAL] parsed type=%s from=%s to=%s callID=%s",
		req.Type,
		req.From,
		req.To,
		req.CallID,
	)

	conn, ok := h.m.Conn(req.To)
	if !ok {

		log.Printf(
			"[SIGNAL] peer not found to=%s",
			req.To,
		)

		wsx.Failed(c)
		return
	}

	log.Printf(
		"[SIGNAL] forward callID=%s",
		req.CallID,
	)

	conn.Do(func(ctx *cira.Context) {
		err := ctx.Push("signal.rtc.forward", req)
		if err != nil {
			log.Printf(
				"[SIGNAL] push failed: %v",
				err,
			)
			return
		}

	})
}

// Register handles a client registration, binding the connection to an ID.
func (h *WSHandler) Register(c *cira.Context) {
	var req RegisterReq
	if !wsx.BindJSON(c, &req) {
		return
	}

	c.Conn.OnClose(func() {
		h.m.Remove(req.ID)
	})
	h.m.Insert(req.ID, c.Conn)

	wsx.OK(c)
}

// UnRegister handles a client unregistration, removing the connection binding.
func (h *WSHandler) UnRegister(c *cira.Context) {
	var req UnRegisterReq
	if !wsx.BindJSON(c, &req) {
		return
	}

	h.m.Remove(req.ID)

	wsx.OK(c)
}
