package meeting

import (
	"github.com/AtoriUzawa/vlink-server/internal/transport/httpx"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the meeting module.
type Handler struct {
	s *Service
}

// NewHandler creates a new Handler with the given Service.
func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

// List handles the meeting room list HTTP endpoint.
func (h *Handler) List(c *gin.Context) {
	httpx.OkWithData(c, h.s.List())
}
