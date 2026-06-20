package p2p

import (
	"sync"

	"github.com/AtoriUzawa/vlink-server/pkg/skiplist"
)

// Manager manages the peer registry with O(log n) lookups via a skip list.
type Manager struct {
	peers map[string]*Peer
	sl    *skiplist.SkipList[*Peer]

	mu sync.RWMutex
}

// NewManager creates a new Manager.
func NewManager() *Manager {
	less := func(a, b *Peer) bool {
		return a.ID < b.ID
	}
	return &Manager{
		peers: make(map[string]*Peer, 0),
		sl:    skiplist.New(less),
	}
}

// Peer returns the peer with the given ID and whether it exists.
func (m *Manager) Peer(id string) (*Peer, bool) {
	p, ok := m.peers[id]
	return p, ok
}

// Register adds a peer to the manager.
func (m *Manager) Register(p *Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.peers[p.ID] = p
	m.sl.Insert(p)
}

// Unregister removes a peer from the manager by ID.
func (m *Manager) Unregister(pid string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	peer, ok := m.peers[pid]
	if !ok {
		return
	}
	delete(m.peers, pid)

	m.sl.Delete(peer)
}

// ListByCursor returns a paginated list of peers starting from the given cursor.
func (m *Manager) ListByCursor(cursor string, limit int) ([]*Peer, string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]*Peer, 0, limit)

	var it *skiplist.Iterator[*Peer]
	if cursor == "" || cursor == "-1" {
		it = m.sl.Iterator()
	} else {
		it = m.sl.SeekAfter(&Peer{ID: cursor})
	}

	for it.Valid() && len(res) < limit {
		res = append(res, it.Value())
		it.Next()
	}

	nextCursor := ""
	if it.Valid() {
		nextCursor = it.Value().ID
	}

	return res, nextCursor
}
