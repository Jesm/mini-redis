package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Store struct {
	values sync.Map
	locks  sync.Map
	timers sync.Map
}

type UnlockCallback func()

func (store *Store) LockKey(key string) UnlockCallback {
	actual, _ := store.locks.LoadOrStore(key, new(sync.Mutex))

	mutex := actual.(*sync.Mutex)
	mutex.Lock()

	return func() {
		mutex.Unlock()
	}
}

type Value interface{}

func (store *Store) Set(key string, value Value) (bool, error) {
	return store.SetEx(key, value, -1)
}

func (store *Store) SetEx(key string, value Value, seconds int) (bool, error) {
	unlock := store.LockKey(key)
	defer unlock()

	store.clearTtlTimer(key)

	store.values.Store(key, value)

	if seconds > -1 {
		store.setTtlTimer(key, seconds)
	}

	return true, nil
}

func (store *Store) setTtlTimer(key string, seconds int) {
	duration := time.Second * time.Duration(seconds)
	timer := time.AfterFunc(duration, func() {
		store.del(key)
	})

	store.timers.Store(key, timer)
}

func (store *Store) clearTtlTimer(key string) {
	if actual, ok := store.timers.Load(key); ok {
		timer := actual.(*time.Timer)
		timer.Stop()

		store.timers.Delete(key)
	}
}

func (store *Store) Get(key string) (string, bool, error) {
	unlock := store.LockKey(key)
	defer unlock()

	if actual, ok := store.values.Load(key); ok {
		switch typed := actual.(type) {
		case string:
			return typed, true, nil
		case int:
			return strconv.Itoa(typed), true, nil
		default:
			return "", false, fmt.Errorf("miniredis: cant return %v of type %T as string", typed, typed)
		}
	}

	return "", false, nil
}

func (store *Store) Del(keys ...string) int {
	count := 0
	for _, key := range keys {
		if store.del(key) {
			count++
		}
	}

	return count
}

func (store *Store) del(key string) bool {
	unlock := store.LockKey(key)
	defer unlock()

	if _, ok := store.values.Load(key); ok {
		store.clearTtlTimer(key)
		store.values.Delete(key)

		return true
	}

	return false
}

func (store *Store) DbSize() int {
	count := 0
	store.values.Range(func(_, _ interface{}) bool {
		count++
		return true
	})

	return count
}

func (store *Store) Incr(key string) (int, error) {
	unlock := store.LockKey(key)
	defer unlock()

	actual, _ := store.values.LoadOrStore(key, 0)

	var num int
	switch typed := actual.(type) {
	case int:
		num = typed
	case string:
		value, err := strconv.Atoi(typed)
		if err != nil {
			return 0, fmt.Errorf("miniredis: conversion of %q to integer failed with message %q", typed, err)
		}

		num = value
	default:
		return 0, fmt.Errorf("miniredis: cant convert value %q to integer", typed)
	}

	num++
	store.values.Store(key, num)

	return num, nil
}

func (store *Store) ZAdd(key string, sets ...SortedSetItem) (int, error) {
	unlock := store.LockKey(key)
	defer unlock()

	actual, _ := store.values.LoadOrStore(key, MakeSortedSet())

	var sortedSet *SortedSet
	switch typed := actual.(type) {
	case *SortedSet:
		sortedSet = typed
	default:
		return 0, fmt.Errorf("miniredis: key %q value is not a sorted set: %q", key, typed)
	}

	count := 0
	for _, set := range sets {
		if sortedSet.Set(set.Score, set.Member) {
			count++
		}
	}

	return count, nil
}

func (store *Store) ZCard(key string) (int, error) {
	unlock := store.LockKey(key)
	defer unlock()

	if actual, ok := store.values.Load(key); ok {
		switch typed := actual.(type) {
		case *SortedSet:
			return typed.Len(), nil
		default:
			return 0, fmt.Errorf("miniredis: key %q value is not a sorted set: %q", key, typed)
		}
	}

	return 0, nil
}

func (store *Store) ZRank(key, member string) (int, bool, error) {
	unlock := store.LockKey(key)
	defer unlock()

	if actual, ok := store.values.Load(key); ok {
		switch typed := actual.(type) {
		case *SortedSet:
			index, ok := typed.Position(member)
			return index, ok, nil
		default:
			return 0, false, fmt.Errorf("miniredis: key %q value is not a sorted set: %q", key, typed)
		}
	}

	return 0, false, nil
}

func (store *Store) ZRange(key string, start, stop int) ([]SortedSetItem, error) {
	unlock := store.LockKey(key)
	defer unlock()

	if actual, ok := store.values.Load(key); ok {
		switch typed := actual.(type) {
		case *SortedSet:
			return typed.Slice(start, stop), nil
		default:
			return nil, fmt.Errorf("miniredis: key %q value is not a sorted set: %q", key, typed)
		}
	}

	return []SortedSetItem{}, nil
}
