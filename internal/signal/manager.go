package signal

import (
	"sync"

	"github.com/AtoriUzawa/cira"
)

// Manager manages WebSocket connections by ID.
type Manager struct {
	conns map[string]*cira.Conn

	mu sync.RWMutex
}

// NewManager creates a new Manager.
func NewManager() *Manager {
	return &Manager{
		conns: make(map[string]*cira.Conn, 0),
	}
}

// Conn returns the WebSocket connection for the given ID and whether it exists.
func (m *Manager) Conn(id string) (*cira.Conn, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, ok := m.conns[id]
	return conn, ok
}

// Insert stores a WebSocket connection associated with the given ID.
func (m *Manager) Insert(id string, conn *cira.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.conns[id] = conn
}

// Remove deletes the WebSocket connection associated with the given ID.
func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.conns, id)
}
