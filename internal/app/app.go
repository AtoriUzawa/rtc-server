// Package app
package app

import (
	"github.com/AtoriUzawa/cira"
	"github.com/AtoriUzawa/vlink-server/internal/live"
	"github.com/AtoriUzawa/vlink-server/internal/meeting"
	"github.com/AtoriUzawa/vlink-server/internal/p2p"
	"github.com/AtoriUzawa/vlink-server/internal/signal"
	"github.com/gin-gonic/gin"
)

// App represents the top-level application, composing all modules and their engines.
type App struct {
	http *gin.Engine
	ws   *cira.Engine
}

// New creates a new App, initializing all modules and registering their routes.
func New() *App {
	// init modules
	signalModule := signal.New()
	p2pModule := p2p.New(signalModule.Manager)
	liveModule := live.New(signalModule.Manager)
	meetingModule := meeting.New(signalModule.Manager)

	// init http
	http := gin.New()
	api := http.Group("/api")

	p2pModule.RegisterHTTP(api)
	liveModule.RegisterHTTP(api)
	meetingModule.RegisterHTTP(api)

	// init ws
	wsEngine := cira.New()

	signalModule.RegisterWS(wsEngine.RouterGroup)
	p2pModule.RegisterWS(wsEngine.RouterGroup)
	liveModule.RegisterWS(wsEngine.RouterGroup)
	meetingModule.RegisterWS(wsEngine.RouterGroup)

	return &App{
		http: http,
		ws:   wsEngine,
	}
}

// RunHTTP starts the HTTP server on the given address.
func (a *App) RunHTTP(addr string) error {
	return a.http.Run(addr)
}

// RunWS starts the WebSocket server on the given address.
func (a *App) RunWS(addr string) error {
	return a.ws.Run(addr)
}
