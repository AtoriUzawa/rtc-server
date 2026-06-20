package meeting

import (
	"sync"
)

// Manager manages meeting rooms with thread-safe operations.
type Manager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

// NewManager creates a new Manager.
func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room, 0),
	}
}

// Insert adds a room to the manager.
func (m *Manager) Insert(r *Room) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.rooms[r.ID] = r
}

// Delete removes a room from the manager by ID.
func (m *Manager) Delete(rid string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.rooms, rid)
}

// Room returns the room with the given ID and whether it exists.
func (m *Manager) Room(rid string) (*Room, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, ok := m.rooms[rid]

	if !ok {
		return nil, false
	}

	return r, true
}

// List returns all rooms in the manager.
func (m *Manager) List() map[string]*Room {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.rooms
}
