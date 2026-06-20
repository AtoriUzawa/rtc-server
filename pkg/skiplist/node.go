package skiplist

type node[T any] struct {
	value   T
	forward []*node[T]
}

func newNode[T any](level int) *node[T] {
	return &node[T]{
		forward: make([]*node[T], level),
	}
}
