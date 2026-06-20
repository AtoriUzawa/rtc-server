// Package heap implements a generic binary heap (priority queue)
//
// Features:
// 1. Supports generic type T
// 2. Priority defined by the less function
// 3. Default implementation uses 0-based complete binary tree
//
// Priority rules:
// less(a, b) == true  => a has higher priority than b
//
// Example:
//
// MinHeap:
//
//	less = func(a, b T) bool { return a < b }
//
// MaxHeap:
//
//	less = func(a, b T) bool { return a > b }
package heap

// Heap generic binary heap structure
type Heap[T any] struct {
	data []T               // array storing heap elements (sequential representation of complete binary tree)
	less func(a, b T) bool // priority comparison function
}

// NewHeap creates an empty heap
func NewHeap[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		data: make([]T, 0),
		less: less,
	}
}

// Len returns the number of elements in the heap
func (h *Heap[T]) Len() int {
	return len(h.data)
}

// Push inserts an element:
// 1. Append to the end of the array
// 2. Sift up to maintain heap property
func (h *Heap[T]) Push(x T) {
	h.data = append(h.data, x)
	h.up(len(h.data) - 1)
}

// Pop removes and returns the top element (highest/lowest priority determined by less):
// 1. Take the root
// 2. Replace the root with the last element
// 3. Remove the last element
// 4. Sift down to maintain heap property
func (h *Heap[T]) Pop() T {
	n := len(h.data)

	x := h.data[0]
	last := h.data[n-1]

	h.data = h.data[:n-1]

	if n > 1 {
		h.data[0] = last
		h.down(0)
	}

	return x
}

// Fix restores the heap property at index i
// Used after externally modifying an element to restore heap invariant
func (h *Heap[T]) Fix(i int) {
	if i > 0 && h.less(h.data[i], h.data[(i-1)/2]) {
		h.up(i)
	} else {
		h.down(i)
	}
}

// swap swaps elements at two positions
func (h *Heap[T]) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

// up sift-up operation:
// swap when child has higher priority than parent
func (h *Heap[T]) up(i int) {
	for i > 0 {
		parent := (i - 1) / 2

		// If the parent has higher priority, the heap property is satisfied
		if h.less(h.data[parent], h.data[i]) {
			break
		}

		h.swap(parent, i)
		i = parent
	}
}

// down sift-down operation:
// ensure parent priority is not lower than children
func (h *Heap[T]) down(i int) {
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

		// current node already satisfies heap property
		if h.less(h.data[i], h.data[priority]) {
			break
		}

		h.swap(priority, i)
		i = priority
	}
}
