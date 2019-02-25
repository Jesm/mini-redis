package main

import (
	"fmt"
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
			return string(typed), true, nil
		default:
			return "", false, fmt.Errorf("miniredis: cant return %q as string", typed)
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
