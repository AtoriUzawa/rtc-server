package meeting

import (
	"github.com/AtoriUzawa/vlink-server/internal/signal"
	"github.com/AtoriUzawa/cira"
	"github.com/gin-gonic/gin"
)

// Module bundles the meeting HTTP handler and WebSocket handler.
type Module struct {
	Handler   *Handler
	WSHandler *WSHandler
}

// New creates a new meeting module with its manager, service, and handlers.
func New(signalManager *signal.Manager) *Module {
	m := NewManager()
	s := NewService(m, signalManager)

	module := &Module{
		Handler:   NewHandler(s),
		WSHandler: NewWSHandler(s),
	}

	return module
}

// RegisterHTTP registers HTTP routes for the meeting module.
func (m *Module) RegisterHTTP(r *gin.RouterGroup) {
	RegisterRouter(r, m.Handler)
}

// RegisterWS registers WebSocket routes for the meeting module.
func (m *Module) RegisterWS(r *cira.RouterGroup) {
	RegisterWSRouter(r, m.WSHandler)
}
