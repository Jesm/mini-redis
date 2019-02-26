package main

import (
	"reflect"
	"testing"
)

func TestMakeSortedSet(t *testing.T) {
	t.Run("make empty sorted set", func(t *testing.T) {
		sortedSet := MakeSortedSet()

		expected := &SortedSet{
			[]SortedSetItem{},
			map[string]int{},
		}
		assertInterface(t, expected, sortedSet)
	})
}

func TestSortSetSet(t *testing.T) {
	sortedSet := MakeSortedSet()

	t.Run("insert score into empty set", func(t *testing.T) {
		if inserted := sortedSet.Set(4, "foo"); !inserted {
			t.Errorf("expected return value to be true")
		}

		expected := []SortedSetItem{{4, "foo"}}
		assertInterface(t, expected, sortedSet.Slice(0, 0))
	})

	t.Run("insert score into non empty set", func(t *testing.T) {
		if inserted := sortedSet.Set(2, "bar"); !inserted {
			t.Errorf("expected return value to be true")
		}

		expected := []SortedSetItem{{2, "bar"}, {4, "foo"}}
		assertInterface(t, expected, sortedSet.Slice(0, 1))
	})

	t.Run("insert score into set with key already present", func(t *testing.T) {
		if inserted := sortedSet.Set(1, "foo"); inserted {
			t.Errorf("expected return value to be false")
		}

		expected := []SortedSetItem{{1, "foo"}, {2, "bar"}}
		assertInterface(t, expected, sortedSet.Slice(0, 1))
	})
}

func TestSortSetLen(t *testing.T) {
	sortedSet := MakeSortedSet()

	t.Run("get cardinality of empty set", func(t *testing.T) {
		assertInterface(t, 0, sortedSet.Len())
	})

	t.Run("get cardinality of set with one element", func(t *testing.T) {
		sortedSet.Set(1, "foo")
		assertInterface(t, 1, sortedSet.Len())
	})

	t.Run("get cardinality of set with many elements", func(t *testing.T) {
		sortedSet.Set(2, "bar")
		sortedSet.Set(3, "xyz")
		assertInterface(t, 3, sortedSet.Len())
	})
}

func TestSortSetPosition(t *testing.T) {
	sortedSet := MakeSortedSet()
	sortedSet.Set(5, "five")
	sortedSet.Set(2, "two")

	t.Run("get position of valid members", func(t *testing.T) {
		if index, ok := sortedSet.Position("two"); ok {
			assertInterface(t, 0, index)
		} else {
			t.Errorf("expected value to be true")
		}

		if index, ok := sortedSet.Position("five"); ok {
			assertInterface(t, 1, index)
		} else {
			t.Errorf("expected value to be true")
		}
	})

	t.Run("get position of invalid member", func(t *testing.T) {
		if _, ok := sortedSet.Position("four"); ok {
			t.Errorf("expected value to be false")
		}
	})
}

type sortedSetSliceTestAux struct {
	start, stop int
	expected    []SortedSetItem
}

func TestSortSetSlice(t *testing.T) {
	sortedSet := MakeSortedSet()
	sortedSet.Set(5, "five")
	sortedSet.Set(2, "two")
	sortedSet.Set(3, "three")
	sortedSet.Set(1, "one")

	t.Run("test multiple ranges", func(t *testing.T) {
		tests := []sortedSetSliceTestAux{
			{0, 3, []SortedSetItem{{1, "one"}, {2, "two"}, {3, "three"}, {5, "five"}}},
			{0, 4, []SortedSetItem{{1, "one"}, {2, "two"}, {3, "three"}, {5, "five"}}},
			{-4, -1, []SortedSetItem{{1, "one"}, {2, "two"}, {3, "three"}, {5, "five"}}},
			{2, 3, []SortedSetItem{{3, "three"}, {5, "five"}}},
			{1, 1, []SortedSetItem{{2, "two"}}},
			{4, 4, []SortedSetItem{}},
			{3, 2, []SortedSetItem{}},
			{5, 2, []SortedSetItem{}},
		}

		for _, te := range tests {
			assertInterface(t, te.expected, sortedSet.Slice(te.start, te.stop))
		}
	})
}

func assertInterface(t *testing.T, expected, got interface{}) {
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
