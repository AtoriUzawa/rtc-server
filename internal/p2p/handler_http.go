package p2p

import (
	"github.com/AtoriUzawa/vlink-server/internal/transport/httpx"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the P2P module.
type Handler struct {
	s *Service
}

// NewHandler creates a new Handler with the given Service.
func NewHandler(s *Service) *Handler {
	return &Handler{
		s: s,
	}
}

// List handles the peer list HTTP endpoint.
func (h *Handler) List(c *gin.Context) {
	var req ListReq
	if !httpx.BindJSON(c, &req) {
		return
	}

	resp := h.s.ListByCursor(&req)
	httpx.OkWithData(c, resp)
}
