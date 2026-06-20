// Package meeting
package meeting

import (
	"sync"
)

// Room represents a meeting room with a host and members.
type Room struct {
	ID      string
	HostID  string
	members map[string]*RoomMember

	mu sync.RWMutex
}

// NewRoom creates a new Room with the given ID.
func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		HostID:  id,
		members: make(map[string]*RoomMember, 0),
	}
}

// Members returns a copy-safe reference to the room's member map.
func (r *Room) Members() map[string]*RoomMember {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.members
}

// Join adds a member to the room.
func (r *Room) Join(m *RoomMember) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.members[m.ID] = m
}

// Leave removes a member from the room by ID.
func (r *Room) Leave(mid string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.members, mid)
}

// Close clears all members from the room.
func (r *Room) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	clear(r.members)
}

// ToDTO converts the room to a RoomDTO for API responses.
func (r *Room) ToDTO() *RoomDTO {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dto := make(map[string]*RoomMemberDTO, len(r.members))
	for k, v := range r.members {
		dto[k] = v.ToDTO()
	}

	return &RoomDTO{
		ID:      r.ID,
		HostID:  r.HostID,
		Members: dto,
	}
}

// RoomMember represents a participant in a meeting room with a role.
type RoomMember struct {
	ID   string
	Role MemberRole
}

// ToDTO converts the room member to a RoomMemberDTO for API responses.
func (m *RoomMember) ToDTO() *RoomMemberDTO {
	return &RoomMemberDTO{
		ID:   m.ID,
		Role: string(m.Role),
	}
}

// MemberRole defines the role of a meeting member.
type MemberRole string

const (
	// Host represents the host role in a meeting.
	Host MemberRole = "host"
	// Member represents a regular participant role in a meeting.
	Member MemberRole = "member"
)
