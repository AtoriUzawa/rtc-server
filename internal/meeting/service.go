package meeting

import (
	"errors"

	"github.com/AtoriUzawa/vlink-server/internal/signal"
	"github.com/AtoriUzawa/cira"
)

// Service provides business logic for meeting operations.
type Service struct {
	m             *Manager
	signalManager *signal.Manager
}

var (
	// ErrMeetingExist is returned when attempting to create a room that already exists.
	ErrMeetingExist = errors.New("the live stream already exist")
	// ErrMeetingNotExist is returned when attempting to access a room that does not exist.
	ErrMeetingNotExist = errors.New("the live stream does not exist")
)

// NewService creates a new Service with the given manager and signal manager.
func NewService(m *Manager, signalManager *signal.Manager) *Service {
	return &Service{
		m:             m,
		signalManager: signalManager,
	}
}

// List returns all meeting rooms as DTOs.
func (s *Service) List() map[string]*RoomDTO {
	mp := s.m.List()

	dto := make(map[string]*RoomDTO, len(mp))
	for k, v := range mp {
		dto[k] = v.ToDTO()
	}

	return dto
}

// Create creates a new meeting room and registers the host connection.
func (s *Service) Create(req *CreateReq, conn *cira.Conn) error {
	_, ok := s.m.Room(req.ID)
	if ok {
		return ErrMeetingExist
	}

	r := NewRoom(req.ID)
	conn.OnClose(func() {
		s.m.Delete(r.ID)
		r.Close()
	})

	s.m.Insert(r)

	r.Join(&RoomMember{r.ID, Host})

	return nil
}

// Join adds a participant to a meeting room.
func (s *Service) Join(req *JoinReq, conn *cira.Conn) error {
	r, ok := s.m.Room(req.RoomID)
	if !ok {
		return ErrMeetingNotExist
	}

	conn.OnClose(func() {
		r.Leave(req.ID)
	})
	r.Join(&RoomMember{req.ID, Member})
	s.broadcastUpdate(r)

	return nil
}

// Leave removes a participant from a meeting room and closes it if the host leaves.
func (s *Service) Leave(req *LeaveReq) error {
	r, ok := s.m.Room(req.RoomID)
	if !ok {
		return ErrMeetingNotExist
	}

	if req.ID == r.HostID {
		s.m.Delete(r.ID)
		r.Leave(req.ID)
		s.broadcastUpdate(r)
		r.Close()
		return nil
	}

	r.Leave(req.ID)
	s.broadcastUpdate(r)

	return nil
}

func (s *Service) broadcastUpdate(r *Room) {
	members := r.Members()
	dto := make(map[string]*RoomMemberDTO, len(members))

	for k, v := range members {
		dto[k] = v.ToDTO()
	}

	for _, member := range members {
		conn, ok := s.signalManager.Conn(member.ID)
		if ok {
			conn.Do(func(c *cira.Context) {
				_ = c.Push("meeting.update", dto)
			})
		}
	}
}
