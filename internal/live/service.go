package live

import (
	"errors"

	"github.com/AtoriUzawa/vlink-server/internal/signal"
	"github.com/AtoriUzawa/cira"
)

// Service provides business logic for live streaming operations.
type Service struct {
	m             *Manager
	signalManager *signal.Manager
}

var (
	// ErrLiveExist is returned when attempting to create a room that already exists.
	ErrLiveExist = errors.New("the live stream already exist")
	// ErrLiveNotExist is returned when attempting to access a room that does not exist.
	ErrLiveNotExist = errors.New("the live stream does not exist")
)

// NewService creates a new Service with the given manager and signal manager.
func NewService(m *Manager, signalManager *signal.Manager) *Service {
	return &Service{
		m:             m,
		signalManager: signalManager,
	}
}

// ListByCursor returns a paginated list of rooms as DTOs.
func (s *Service) ListByCursor(req *ListReq) *ListResp {
	req.Limit = max(1, req.Limit)
	req.Limit = min(10, req.Limit)
	rs, nextCursor := s.m.ListByCursor(req.Cursor, req.Limit)
	dtoRs := make([]*RoomDTO, len(rs))

	for i, r := range rs {
		dtoRs[i] = r.ToDTO()
	}

	return &ListResp{
		Rooms:      dtoRs,
		NextCursor: nextCursor,
	}
}

// Create creates a new live room and registers the owner connection.
func (s *Service) Create(req *CreateReq, conn *cira.Conn) error {
	_, ok := s.m.Room(req.ID)
	if ok {
		return ErrLiveExist
	}

	r := NewRoom(req.ID, req.Title, conn)
	conn.OnClose(func() {
		s.m.Delete(r.ID)
		r.Close()
	})

	s.m.Insert(r)

	return nil
}

// Join adds a participant to a live room.
func (s *Service) Join(req *JoinReq, conn *cira.Conn) error {
	// refresh skiplist
	err := s.m.Update(req.RoomID, func(r *Room) {
		r.Join(&RoomMember{
			ID:   req.ID,
			Role: RoleViewer,
		})

		conn.OnClose(func() {
			r.Leave(req.ID)
		})

		s.broadcastUpdate(r)
	})
	if err != nil {
		return ErrLiveNotExist
	}

	return nil
}

// Leave removes a participant from a live room and closes it if the owner leaves.
func (s *Service) Leave(req *LeaveReq) error {
	r, ok := s.m.Room(req.RoomID)
	if !ok {
		return ErrLiveNotExist
	}

	err := s.m.Update(req.RoomID, func(r *Room) {
		r.Leave(req.ID)
		s.broadcastUpdate(r)
	})
	if err != nil {
		return ErrLiveNotExist
	}

	if r.ID == req.ID {
		s.m.Delete(r.ID)
		r.Close()
	}

	return nil
}

func (s *Service) broadcastUpdate(r *Room) {
	members := r.Members()
	dto := make(map[string]*RoomMemberDTO, len(members))
	for i, m := range members {
		dto[i] = m.ToDTO()
	}

	for _, member := range members {
		conn, ok := s.signalManager.Conn(member.ID)
		if ok {
			conn.Do(func(c *cira.Context) {
				_ = c.Push("live.update", dto)
			})
		}
	}
}
