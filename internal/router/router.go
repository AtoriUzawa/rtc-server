// Package router
package router

import (
	"github.com/AtoriUzawa/vlink-server/internal/live"
	"github.com/AtoriUzawa/vlink-server/internal/p2p"
	"github.com/AtoriUzawa/cira"
	"github.com/gin-gonic/gin"
)

// Handlers aggregates module HTTP handlers for registration.
type Handlers struct {
	P2P  *p2p.Handler
	Room *live.Handler
}
// WSHandler aggregates module WebSocket handlers for registration.
type WSHandler struct {
	P2P  *p2p.WSHandler
	Room *live.WSHandler
}

// Register registers HTTP routes with the provided gin engine and handlers.
func Register(r *gin.Engine, h Handlers) {
	api := r.Group("/api")

	p2p.RegisterRouter(api, h.P2P)
	live.RegisterRouter(api, h.Room)
}

// WSRegister registers WebSocket routes with the provided ws engine and handlers.
func WSRegister(r *cira.Engine, h WSHandler) {
	event := r.Group("rtc")
	p2p.RegisterWSRouter(event, h.P2P)
	live.RegisterWSRouter(event, h.Room)
}
