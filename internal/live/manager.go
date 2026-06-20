package live

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/AtoriUzawa/vlink-server/pkg/skiplist"
)

// Manager manages live rooms with a skip list for heat-based ordering.
type Manager struct {
	rooms map[string]*Room
	sl    *skiplist.SkipList[*RoomItem]

	mu sync.RWMutex
}

// NewManager creates a new Manager.
func NewManager() *Manager {
	less := func(a, b *RoomItem) bool {
		if a.Heat != b.Heat {
			return a.Heat > b.Heat
		}

		return a.ID < b.ID
	}
	return &Manager{
		rooms: make(map[string]*Room, 0),
		sl:    skiplist.New(less),
	}
}

// Insert adds a room to the manager.
func (m *Manager) Insert(r *Room) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.rooms[r.ID] = r
	m.sl.Insert(r.ToItem())
}

// Delete removes a room from the manager by ID.
func (m *Manager) Delete(rid string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, ok := m.rooms[rid]
	if !ok {
		return
	}

	delete(m.rooms, rid)
	m.sl.Delete(r.ToItem())
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

// Update can be extended to worker processing
// Do not call Manager methods within it to avoid lock contention
func (m *Manager) Update(rid string, fn func(r *Room)) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, ok := m.rooms[rid]
	if !ok {
		return errors.New("ErrRoomNotFound")
	}
	m.sl.Delete(r.ToItem())

	fn(r)

	m.sl.Insert(r.ToItem())

	return nil
}

// ListByCursor returns a paginated list of rooms sorted by heat, starting from the given cursor.
func (m *Manager) ListByCursor(cursor string, limit int) ([]*Room, string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nextCursor := ""
	res := make([]*Room, 0, limit)

	var it *skiplist.Iterator[*RoomItem]

	item, err := DecodeCursor(cursor)
	if err != nil {
		it = m.sl.Iterator()
	} else {
		it = m.sl.SeekAfter(item)
	}

	for it.Valid() && len(res) < limit {
		res = append(res, m.rooms[it.Value().ID])
		it.Next()
	}

	if it.Valid() {
		nextCursor = EncodeCursor(it.Value())
	}

	return res, nextCursor
}

// EncodeCursor encodes a RoomItem into a cursor string for pagination.
func EncodeCursor(item *RoomItem) string {
	return item.ID + "|" + strconv.Itoa(item.Heat)
}

// DecodeCursor decodes a cursor string back into a RoomItem.
func DecodeCursor(cursor string) (*RoomItem, error) {
	split := strings.Split(cursor, "|")
	if len(split) != 2 {
		return nil, errors.New("invalid cursor")
	}

	l, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, err
	}

	return &RoomItem{
		ID:   split[0],
		Heat: l * 10,
	}, nil
}
