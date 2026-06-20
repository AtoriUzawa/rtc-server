package signal

import (
	"github.com/AtoriUzawa/cira"
)

// Module bundles the signal manager and WebSocket handler.
type Module struct {
	Manager   *Manager
	WSHandler *WSHandler
}

// New creates a new signal module with its manager and WebSocket handler.
func New() *Module {
	m := NewManager()

	module := &Module{
		Manager:   m,
		WSHandler: NewWSHandler(m),
	}

	return module
}

// RegisterWS registers WebSocket routes for the signal module.
func (m *Module) RegisterWS(r *cira.RouterGroup) {
	RegisterWSRouter(r, m.WSHandler)
}
