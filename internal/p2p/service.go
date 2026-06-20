package p2p

import (
	"log"
	"time"

	"github.com/AtoriUzawa/vlink-server/internal/signal"
	"github.com/AtoriUzawa/cira"
)

// Service provides business logic for P2P operations.
type Service struct {
	m             *Manager
	signalManager *signal.Manager
}

// NewService creates a new Service with the given manager and signal manager.
func NewService(m *Manager, signalManager *signal.Manager) *Service {
	return &Service{
		m:             m,
		signalManager: signalManager,
	}
}

// ListByCursor returns a paginated list of peers as DTOs.
func (s *Service) ListByCursor(req *ListReq) *ListResp {
	req.Limit = max(1, req.Limit)
	req.Limit = min(10, req.Limit)
	ps, nextCursor := s.m.ListByCursor(req.Cursor, req.Limit)
	dtoPs := make([]*PeerDTO, len(ps))
	for i, p := range ps {
		dtoPs[i] = &PeerDTO{
			ID:       p.ID,
			Nickname: "Atori",
		}
	}
	return &ListResp{
		Peers:      dtoPs,
		NextCursor: nextCursor,
	}
}

// Call initiates a call to a peer via signaling and returns the peer's response.
func (s *Service) Call(c *cira.Context, req *CallReq) {
	log.Printf(
		"[P2P][CALL] from=%s to=%s callID=%s",
		req.From,
		req.To,
		req.CallID,
	)

	peer, ok := s.m.Peer(req.To)
	if !ok {
		log.Printf(
			"[P2P][CALL] target offline from=%s to=%s",
			req.From,
			req.To,
		)

		c.Resp(map[string]string{
			"status": "offline",
		})
		return
	}

	log.Printf(
		"[P2P][CALL] target found id=%s",
		peer.ID,
	)

	conn, ok := s.signalManager.Conn(peer.ID)
	if !ok {
		log.Printf(
			"[P2P][CALL] target ws disconnected id=%s",
			peer.ID,
		)

		c.Resp(map[string]string{
			"status": "offline",
		})
		return
	}

	var (
		resp    map[string]string
		callErr error
	)

	log.Printf(
		"[P2P][CALL] forwarding request -> %s",
		peer.ID,
	)

	conn.Do(func(ctx *cira.Context) {
		ctx.Timeout = 30 * time.Second

		callErr = ctx.Call(
			"p2p.call",
			req,
			&resp,
		)
	})

	if callErr != nil {
		log.Printf(
			"[P2P][CALL] timeout from=%s to=%s err=%v",
			req.From,
			req.To,
			callErr,
		)

		c.Resp(map[string]string{
			"status": "timeout",
		})
		return
	}

	log.Printf(
		"[P2P][CALL] response from=%s to=%s status=%v",
		req.To,
		req.From,
		resp,
	)

	c.Resp(resp)
}

// Hangup sends a hangup notification to the target peer.
func (s *Service) Hangup(req *HangupReq) {
	log.Printf(
		"[P2P][HANGUP] from=%s to=%s",
		req.From,
		req.To,
	)

	peer, ok := s.m.Peer(req.To)
	if !ok {
		log.Printf(
			"[P2P][HANGUP] target offline to=%s",
			req.To,
		)
		return
	}

	conn, ok := s.signalManager.Conn(peer.ID)
	if !ok {
		log.Printf(
			"[P2P][HANGUP] target ws disconnected to=%s",
			req.To,
		)
		return
	}

	conn.Do(func(c *cira.Context) {
		log.Printf(
			"[P2P][HANGUP] forwarding hangup -> %s",
			req.To,
		)

		_ = c.Push(
			"p2p.hangup",
			req,
		)
	})
}
