// Package skiplist provides a generic skip list implementation.
//
// SkipList is an ordered probabilistic data structure that supports
// average O(log n) insertion, deletion and lookup operations.
//
// Ordering is determined by the user-provided less function.
// If neither less(a,b) nor less(b,a) is true,
// a and b are considered equal.
//
// SkipList is NOT safe for concurrent use.
// Callers must synchronize access externally.
package skiplist

import (
	"math/rand/v2"
)

// SkipList is an ordered probabilistic data structure supporting O(log n)
// average insertion, deletion and lookup operations.
type SkipList[T any] struct {
	head     *node[T]
	size     int
	level    int
	maxLevel int
	p        float64

	less func(a, b T) bool
	rand *rand.Rand
}

const (
	// DefaultMaxLevel is the default maximum level for a skip list node.
	DefaultMaxLevel = 32
	// DefaultProbability is the default probability used for random level generation.
	DefaultProbability = 0.25
)

// Option is a functional option for configuring a SkipList.
type Option[T any] func(*SkipList[T])

// WithProbability returns an Option that sets the probability for random level generation.
// The probability must be in the range (0, 1).
func WithProbability[T any](p float64) Option[T] {
	if p <= 0 || p >= 1 {
		panic("skiplist: probability must be in (0,1)")
	}

	return func(sl *SkipList[T]) {
		sl.p = p
	}
}

// WithMaxLevel returns an Option that sets the maximum level for a skip list node.
func WithMaxLevel[T any](maxLevel int) Option[T] {
	if maxLevel <= 0 {
		panic("skiplist: maxLevel must > 0")
	}

	return func(sl *SkipList[T]) {
		sl.maxLevel = maxLevel
	}
}

// New creates a new SkipList ordered by the given less function, with optional configuration.
func New[T any](less func(a, b T) bool, opts ...Option[T]) *SkipList[T] {
	sl := &SkipList[T]{
		level:    1,
		maxLevel: DefaultMaxLevel,
		p:        DefaultProbability,
		less:     less,
		rand:     rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
	}

	for _, opt := range opts {
		opt(sl)
	}

	sl.head = newNode[T](sl.maxLevel)

	return sl
}

// Len returns the number of elements in the skip list.
func (s *SkipList[T]) Len() int {
	return s.size
}

// Search looks up v in the skip list.
//
// It returns the stored value and true if found.
// Otherwise it returns the zero value of T and false.
func (s *SkipList[T]) Search(v T) (T, bool) {
	cur := s.findPath(v, nil)
	cur = cur.forward[0]
	if cur != nil && s.equal(cur.value, v) {
		return cur.value, true
	}

	return zero[T](), false
}

// lowerBound returns the first node whose value is greater than
// or equal to v.
//
// Equivalent to:
//
// first element >= v
func (s *SkipList[T]) lowerBound(v T) *node[T] {
	cur := s.findPath(v, nil)
	return cur.forward[0]
}

// upperBound returns the first node whose value is strictly
// greater than v.
//
// Equivalent to:
//
// first element > v
func (s *SkipList[T]) upperBound(v T) *node[T] {
	cur := s.findPath(v, nil)
	cur = cur.forward[0]
	if cur == nil {
		return nil
	}

	if s.equal(cur.value, v) {
		cur = cur.forward[0]
	}

	return cur
}

// Insert inserts v into the skip list.
//
// Returns false if an equal value already exists.
func (s *SkipList[T]) Insert(v T) bool {
	update := make([]*node[T], s.maxLevel)
	cur := s.findPath(v, update)

	level := s.randLevel()

	candidate := cur.forward[0]
	if candidate != nil && s.equal(candidate.value, v) {
		return false
	}

	if level > s.level {
		for i := s.level; i < level; i++ {
			update[i] = s.head
		}
		s.level = level
	}

	n := newNode[T](level)
	n.value = v

	for i := level - 1; i >= 0; i-- {
		n.forward[i] = update[i].forward[i]
		update[i].forward[i] = n
	}

	s.size++
	return true
}

// Delete removes v from the skip list.
//
// Returns true if v existed and was removed.
func (s *SkipList[T]) Delete(v T) bool {
	update := make([]*node[T], s.maxLevel)
	cur := s.findPath(v, update)

	cur = cur.forward[0]
	if cur != nil && s.equal(cur.value, v) {
		for i := len(cur.forward) - 1; i >= 0; i-- {
			update[i].forward[i] = cur.forward[i]
		}

		for s.level > 1 && s.head.forward[s.level-1] == nil {
			s.level--
		}

		s.size--
		return true
	}

	return false
}

// Iterator returns an iterator positioned at the first element.
//
// If the skip list is empty, the returned iterator is invalid.
func (s *SkipList[T]) Iterator() *Iterator[T] {
	return &Iterator[T]{
		cur: s.head.forward[0],
	}
}

// Seek returns an iterator positioned at the first element
// whose value is greater than or equal to v.
//
// Equivalent to lowerBound(v).
func (s *SkipList[T]) Seek(v T) *Iterator[T] {
	return &Iterator[T]{cur: s.lowerBound(v)}
}

// SeekAfter returns an iterator positioned at the first element
// whose value is strictly greater than v.
//
// Equivalent to upperBound(v).
//
// This method is useful for cursor-based pagination.
func (s *SkipList[T]) SeekAfter(v T) *Iterator[T] {
	return &Iterator[T]{cur: s.upperBound(v)}
}

// Range returns at most limit elements starting from the first
// element greater than or equal to start.
//
// Results are returned in sorted order.
func (s *SkipList[T]) Range(start T, limit int) []T {
	it := s.Seek(start)
	res := make([]T, 0, limit)
	for it.Valid() && len(res) < limit {
		res = append(res, it.Value())
		it.Next()
	}

	return res
}

// RangeByValue returns all elements in the half-open interval:
//
// [start, end)
//
// Results are returned in sorted order.
func (s *SkipList[T]) RangeByValue(start, end T) []T {
	it := s.Seek(start)
	res := make([]T, 0)
	for it.Valid() && s.less(it.Value(), end) {
		res = append(res, it.Value())
		it.Next()
	}

	return res
}

// randLevel generates a random level for a new node.
//
// Higher levels occur with probability p.
func (s *SkipList[T]) randLevel() int {
	level := 1

	for level < s.maxLevel && s.rand.Float64() < s.p {
		level++
	}

	return level
}

// equal reports whether a and b are considered equal
// under the current ordering relation.
func (s *SkipList[T]) equal(a, b T) bool {
	return !s.less(a, b) && !s.less(b, a)
}

// findPath finds the predecessor path of v.
//
// If update is not nil, update[i] stores the predecessor node
// at level i.
//
// The returned node is the last node whose value is less than v.
func (s *SkipList[T]) findPath(v T, update []*node[T]) *node[T] {
	cur := s.head
	for i := s.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && s.less(cur.forward[i].value, v) {
			cur = cur.forward[i]
		}
		if update != nil {
			update[i] = cur
		}
	}
	return cur
}

// zero returns the zero value of T.
func zero[T any]() T {
	var z T
	return z
}
