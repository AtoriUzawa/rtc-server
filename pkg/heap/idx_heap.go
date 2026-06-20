package heap

// IdxHeap is an indexed heap (Indexed Priority Queue)
//
// Features:
// 1. Supports O(log n) Push / Pop / Remove
// 2. Supports O(1) element lookup by key
// 3. Maintains "element -> heap position" mapping via idx map
//
// T: element type
// K: unique element identifier (must be comparable, used as map key)
type IdxHeap[T any, K comparable] struct {
	data []T               // heap array (sequential storage of complete binary tree)
	less func(a, b T) bool // priority function: true means a has higher priority
	idx  map[K]int         // key -> index mapping table
	key  func(T) K         // extract unique identifier K from T
}

// NewIdxHeap creates an Indexed Heap
func NewIdxHeap[T any, K comparable](
	less func(a, b T) bool,
	key func(T) K,
) *IdxHeap[T, K] {
	return &IdxHeap[T, K]{
		data: make([]T, 0),
		less: less,
		idx:  make(map[K]int),
		key:  key,
	}
}

// Len returns the heap size
func (h *IdxHeap[T, K]) Len() int {
	return len(h.data)
}

// Push inserts an element:
// 1. Insert at the end of the array
// 2. Update idx
// 3. Sift up to maintain heap property
func (h *IdxHeap[T, K]) Push(x T) {
	h.data = append(h.data, x)
	i := len(h.data) - 1

	h.idx[h.key(x)] = i

	h.up(i)
}

// Pop removes and returns the top element:
// 1. Take the root
// 2. Replace root with the last element
// 3. Remove last
// 4. Fix idx
// 5. Sift down to maintain heap property
func (h *IdxHeap[T, K]) Pop() T {
	n := len(h.data)
	x := h.data[0]

	last := n - 1

	// Swap the last element to root
	h.swap(0, last)

	// Remove last
	h.data = h.data[:last]
	delete(h.idx, h.key(x))

	// Fix heap
	if last > 0 {
		h.down(0)
	}

	return x
}

// Remove removes any element by key:
// 1. O(1) lookup of index
// 2. Swap to last position
// 3. Remove last
// 4. Fix heap structure
func (h *IdxHeap[T, K]) Remove(key K) T {
	i := h.idx[key]

	last := len(h.data) - 1

	// Swap the target element with the last element
	h.swap(i, last)

	// Element to be removed (now at the end)
	remove := h.data[last]

	// Remove last
	h.data = h.data[:last]
	delete(h.idx, key)

	// Fix the swapped position
	if i < len(h.data) {
		h.fix(i)
	}

	return remove
}

// Slice returns a slice over the specified range
func (h *IdxHeap[T, K]) Slice(i, j int) []T {
	n := len(h.data)
	if i < 0 {
		i = 0
	}
	if j > n {
		j = n
	}
	if i >= j {
		return []T{}
	}
	res := make([]T, j-i)
	copy(res, h.data[i:j])
	return res
}

func (h *IdxHeap[T, K]) Key(k T) (int, bool) {
	idx, bool := h.idx[h.key(k)]
	return idx, bool
}

// fix restores the heap property at a position:
// may sift up or down (chooses one direction)
func (h *IdxHeap[T, K]) fix(i int) {
	if i > 0 && h.less(h.data[i], h.data[(i-1)/2]) {
		h.up(i)
	} else {
		h.down(i)
	}
}

// swap swaps two positions in the heap and updates the idx table
func (h *IdxHeap[T, K]) swap(i, j int) {
	a := h.data[i]
	b := h.data[j]

	h.data[i], h.data[j] = h.data[j], h.data[i]

	// Update index mapping
	h.idx[h.key(a)] = j
	h.idx[h.key(b)] = i
}

// up sift-up operation (maintains heap: child cannot have higher priority than parent)
func (h *IdxHeap[T, K]) up(i int) {
	for i > 0 {
		parent := (i - 1) / 2

		if h.less(h.data[parent], h.data[i]) {
			break
		}

		h.swap(parent, i)
		i = parent
	}
}

// down sift-down operation (maintains heap: parent cannot have lower priority than child)
func (h *IdxHeap[T, K]) down(i int) {
	n := len(h.data)

	for {
		l := i*2 + 1
		if l >= n {
			break
		}

		r := l + 1

		// select the child with higher priority
		priority := l
		if r < n && h.less(h.data[r], h.data[l]) {
			priority = r
		}

		// if the current node already satisfies the heap property, stop
		if h.less(h.data[i], h.data[priority]) {
			break
		}

		h.swap(priority, i)
		i = priority
	}
}
