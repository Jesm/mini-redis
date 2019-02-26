package main

import (
	"sort"
)

type SortedSet struct {
	items []SortedSetItem
	index map[string]int
}

func MakeSortedSet() *SortedSet {
	return &SortedSet{
		make([]SortedSetItem, 0),
		make(map[string]int),
	}
}

type SortedSetItem struct {
	Score  float64
	Member string
}

func (set *SortedSet) Set(score float64, member string) bool {
	item := SortedSetItem{score, member}
	defer set.ensureOrder()

	if index, ok := set.index[member]; ok {
		set.items[index] = item
		return false
	}

	set.items = append(set.items, item)
	return true
}

func (set *SortedSet) ensureOrder() {
	sort.SliceStable(set.items, func(i, j int) bool {
		return set.items[i].Score < set.items[j].Score
	})

	for index, item := range set.items {
		set.index[item.Member] = index
	}
}

func (set *SortedSet) Len() int {
	return len(set.items)
}

func (set *SortedSet) Position(member string) (int, bool) {
	index, ok := set.index[member]
	return index, ok
}

func (set *SortedSet) Slice(start, stop int) []SortedSetItem {
	size := len(set.items)
	if start >= size || start > stop {
		return []SortedSetItem{}
	}

	if start < 0 {
		if start += size; start < 0 {
			return []SortedSetItem{}
		}
	}

	switch {
	case stop >= size:
		stop = size
	case stop < 0:
		stop += size + 1
	default:
		stop++
	}

	return append([]SortedSetItem{}, set.items[start:stop]...)
}
