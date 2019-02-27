package main

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	store := new(Store)
	store.Set("foo", "1")

	setTestAux(t, store, "insert value into non existing key", "bar", "2")
	setTestAux(t, store, "insert value into existing key", "foo", "3")
}

func setTestAux(t *testing.T, store *Store, name, key, value string) {
	t.Run(name, func(t *testing.T) {
		if ok, err := store.Set(key, value); err == nil {
			if !ok {
				t.Errorf("expected true, got false")
			}
		} else {
			t.Errorf("expected nil, got %q", err)
		}

		assertGet(t, store, key, value, true, false)
	})
}

func TestSetEx(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test")
	}

	store := new(Store)
	store.SetEx("foo", "1", 1)

	t.Run("get key before expiration", func(t *testing.T) {
		assertGet(t, store, "foo", "1", true, false)
	})

	t.Run("get key after expiration", func(t *testing.T) {
		ch := make(chan bool)
		dur := time.Second * time.Duration(2)
		time.AfterFunc(dur, func() {
			assertGet(t, store, "foo", "", false, false)
			close(ch)
		})

		<-ch
	})
}

func TestGet(t *testing.T) {
	store := new(Store)
	store.Set("foo", "bar")
	store.Set("float", float64(1))

	t.Run("get valid existing key", func(t *testing.T) {
		assertGet(t, store, "foo", "bar", true, false)
	})
	t.Run("get non existing key", func(t *testing.T) {
		assertGet(t, store, "bar", "", false, false)
	})
	t.Run("get invalid key", func(t *testing.T) {
		assertGet(t, store, "float", "", false, true)
	})
}

func TestDel(t *testing.T) {
	store := new(Store)
	store.Set("foo", "bar")

	t.Run("delete key", func(t *testing.T) {
		store.Del("foo")
		assertGet(t, store, "foo", "", false, false)
	})

	t.Run("delete multiple keys, including non existing one", func(t *testing.T) {
		store.Set("fizz", "buzz")
		store.Set("buzz", "fizz")
		count := store.Del("foo", "fizz", "buzz")

		if expected := 2; count != expected {
			t.Errorf("expected %v, got %v", expected, count)
		}
	})
}

func assertGet(t *testing.T, store *Store, key, expected string, expectedOk, foundErr bool) {
	actual, ok, err := store.Get(key)

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}

	if ok != expectedOk {
		t.Errorf("expected %v, got %v", expectedOk, ok)
	}

	if foundErr && err == nil {
		t.Errorf("expected error, but got nil")
	} else if !foundErr && err != nil {
		t.Errorf("expected nil, got %q", err)
	}
}

func TestDbSize(t *testing.T) {
	store := new(Store)
	store.Set("foo", "br")
	store.Set("fizz", "buzz")
	store.Set("foo", "bar")

	t.Run("get the store size", func(t *testing.T) {
		count := store.DbSize()
		if expected := 2; count != expected {
			t.Errorf("expected %v, got %v", expected, count)
		}
	})
}

func TestIncr(t *testing.T) {
	store := new(Store)
	store.Set("bar", "4")
	store.Set("fizz", "buzz")

	t.Run("increment non existing key", func(t *testing.T) {
		num, _ := store.Incr("foo")
		if expected := 1; num != expected {
			t.Errorf("expected %v, got %v", expected, num)
		}
	})

	t.Run("increment existing key", func(t *testing.T) {
		num, _ := store.Incr("bar")
		if expected := 5; num != expected {
			t.Errorf("expected %v, got %v", expected, num)
		}
	})

	t.Run("increment invalid key", func(t *testing.T) {
		_, err := store.Incr("fizz")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("increment key in multiple goroutines", func(t *testing.T) {
		wg := new(sync.WaitGroup)
		for x := 0; x < 1000; x++ {
			wg.Add(1)

			go func() {
				store.Incr("xyz")
				wg.Done()
			}()
		}

		wg.Wait()
		assertGet(t, store, "xyz", "1000", true, false)
	})
}

func TestZAdd(t *testing.T) {
	store := new(Store)
	store.Set("foo", "bar")

	t.Run("add to non existing key", func(t *testing.T) {
		if count, err := store.ZAdd("fizz", SortedSetItem{8, "eight"}); err == nil {
			if expected := 1; count != expected {
				t.Errorf("expected %v, got %v", expected, count)
			}
		} else {
			t.Errorf("expected nil, got %q", err)
		}
	})

	t.Run("add to key with invalid value", func(t *testing.T) {
		if _, err := store.ZAdd("foo", SortedSetItem{4, "four"}); err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestZCard(t *testing.T) {
	store := new(Store)
	store.Set("foo", "bar")

	t.Run("get cardinality of non existing key", func(t *testing.T) {
		if count, err := store.ZCard("fizz"); err == nil {
			if expected := 0; count != expected {
				t.Errorf("expected %v, got %v", expected, count)
			}
		} else {
			t.Errorf("expected nil, got %q", err)
		}
	})

	t.Run("get cardinality of key with invalid value", func(t *testing.T) {
		if _, err := store.ZCard("foo"); err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestZRank(t *testing.T) {
	store := new(Store)
	store.Set("foo", "bar")

	t.Run("get rank of non existing key", func(t *testing.T) {
		if _, ok, err := store.ZRank("fizz", "member"); err == nil {
			if ok {
				t.Errorf("expected false, got true")
			}
		} else {
			t.Errorf("expected nil, got %q", err)
		}
	})

	t.Run("get rank of key with invalid value", func(t *testing.T) {
		if _, _, err := store.ZRank("foo", "member"); err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestZRange(t *testing.T) {
	store := new(Store)
	store.Set("foo", "bar")

	t.Run("get rank of non existing key", func(t *testing.T) {
		if arr, err := store.ZRange("fizz", 0, 1); err == nil {
			if expected := []SortedSetItem{}; !reflect.DeepEqual(expected, arr) {
				t.Errorf("expected %v, got %v", expected, arr)
			}
		} else {
			t.Errorf("expected nil, got %q", err)
		}
	})

	t.Run("get range of key with invalid value", func(t *testing.T) {
		if _, err := store.ZRange("foo", 0, 2); err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
