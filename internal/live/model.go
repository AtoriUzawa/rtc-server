// Package live
package live

import (
	"sync"

	"github.com/AtoriUzawa/cira"
)

// Room represents a live streaming room with an owner and members.
type Room struct {
	ID      string
	Title   string
	OwnerID string
	members map[string]*RoomMember

	mu sync.RWMutex
}

// NewRoom creates a new Room with the given owner, title, and connection.
func NewRoom(ownerID string, title string, conn *cira.Conn) *Room {
	r := &Room{
		ID:      ownerID,
		Title:   title,
		OwnerID: ownerID,
		members: make(map[string]*RoomMember, 0),
	}

	r.members[ownerID] = &RoomMember{
		ID:   ownerID,
		Role: RoleOwner,
	}

	return r
}

// Members returns a copy-safe reference to the room's member map.
func (r *Room) Members() map[string]*RoomMember {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.members
}

// Heat returns the room's popularity score based on member count.
func (r *Room) Heat() int {
	return len(r.members) * 10
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

// ToItem converts the room to a RoomItem for skip list ordering.
func (r *Room) ToItem() *RoomItem {
	return &RoomItem{
		ID:   r.ID,
		Heat: r.Heat(),
	}
}

// ToDTO converts the room to a RoomDTO for API responses.
func (r *Room) ToDTO() *RoomDTO {
	r.mu.Lock()
	defer r.mu.Unlock()

	return &RoomDTO{
		ID:      r.OwnerID,
		Title:   r.Title,
		OwnerID: r.OwnerID,
		Count:   len(r.members),
	}
}

// RoomMember represents a participant in a live room with a role.
type RoomMember struct {
	ID   string
	Role MemberRole
}

// ToDTO converts the room member to a RoomMemberDTO for API responses.
func (rm *RoomMember) ToDTO() *RoomMemberDTO {
	return &RoomMemberDTO{
		ID:   rm.ID,
		Role: string(rm.Role),
	}
}

// MemberRole defines the role of a room member.
type MemberRole string

const (
	// RoleViewer represents a viewer role in a live room.
	RoleViewer MemberRole = "viewer"
	// RoleOwner represents the owner role in a live room.
	RoleOwner MemberRole = "owner"
)

// RoomItem is a lightweight representation used for skip list ordering by heat.
type RoomItem struct {
	ID   string
	Heat int
}
