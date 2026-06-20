package skiplist

import (
	"testing"
)

func TestInsertAndSearch(t *testing.T) {
	sl := New(func(a, b int) bool {
		return a < b
	})

	tl := []int{1, 3, 5, 2, 9, 8}
	for _, v := range tl {
		sl.Insert(v)
	}

	for _, v := range tl {
		_, ok := sl.Search(v)
		if !ok {
			t.Fatalf("%d should exist", v)
		}

	}

	if _, ok := sl.Search(-1); ok {
		t.Fatal("-1 should not exist")
	}

	if !sl.Delete(1) {
		t.Fatal("1 should be deleted")
	}

	if _, ok := sl.Search(1); ok {
		t.Fatal("1 should not exist after delete")
	}

	if sl.Delete(1) {
		t.Fatal("1 should not be deleted twice")
	}

	slLen := sl.Len()

	if slLen != 5 {
		t.Fatalf("expected len=5, got=%d", slLen)
	}
}

func TestDuplicateInsert(t *testing.T) {
	sl := New(func(a, b int) bool {
		return a < b
	})

	if !sl.Insert(1) {
		t.Fatal("first insert should succeed")
	}

	if sl.Insert(1) {
		t.Fatal("second insert should fail")
	}

	if sl.Len() != 1 {
		t.Fatalf("expected len=1, got=%d", sl.Len())
	}
}

func TestEmptySkipList(t *testing.T) {
	sl := New(func(a, b int) bool {
		return a < b
	})

	if _, ok := sl.Search(1); ok {
		t.Fatal("empty list should not contain anything")
	}

	if sl.Delete(1) {
		t.Fatal("delete on empty list should fail")
	}

	if sl.Len() != 0 {
		t.Fatal("empty list len should be 0")
	}
}

func TestRandomInsert(t *testing.T) {
	sl := New(func(a, b int) bool {
		return a < b
	})

	for i := range 10000 {
		sl.Insert(i)
	}

	for i := range 10000 {
		if _, ok := sl.Search(i); !ok {
			t.Fatalf("%d missing", i)
		}
	}
}

func TestOrder(t *testing.T) {
	sl := New(func(a, b int) bool {
		return a < b
	})

	tl := []int{5, 1, 9, 3, 7}

	for _, v := range tl {
		sl.Insert(v)
	}

	cur := sl.head.forward[0]

	prev := -1

	for cur != nil {
		if cur.value < prev {
			t.Fatal("order broken")
		}

		prev = cur.value
		cur = cur.forward[0]
	}
}

func TestDeleteAll(t *testing.T) {
	sl := New(func(a, b int) bool {
		return a < b
	})

	for i := range 100 {
		sl.Insert(i)
	}

	for i := range 100 {
		if !sl.Delete(i) {
			t.Fatalf("%d delete failed", i)
		}
	}

	if sl.Len() != 0 {
		t.Fatal("len should be 0")
	}
}
