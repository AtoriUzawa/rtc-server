package skiplist

// Iterator iterates over the elements of a SkipList in sorted order.
type Iterator[T any] struct {
	cur *node[T]
}

// Next advances the iterator to the next element.
func (it *Iterator[T]) Next() {
	it.cur = it.cur.forward[0]
}

// Value returns the current element. The iterator must be valid.
func (it *Iterator[T]) Value() T {
	return it.cur.value
}

// Valid reports whether the iterator is positioned at a valid element.
func (it *Iterator[T]) Valid() bool {
	return it.cur != nil
}
