package main

import (
	"reflect"
	"testing"
)

func TestExec(t *testing.T) {
	store := new(Store)
	intr := Interpreter{store}

	t.Run("set key to store", func(t *testing.T) {
		if actual, err := intr.Exec("SET foo bar"); err == nil {
			if expected := true; actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("expected no error, but got %q", err)
		}
	})

	t.Run("get key from store", func(t *testing.T) {
		if actual, err := intr.Exec("GET foo"); err == nil {
			if expected := "bar"; actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("expected no error, but got %q", err)
		}
	})

	t.Run("get the store size", func(t *testing.T) {
		if actual, err := intr.Exec("DBSIZE"); err == nil {
			if expected := 1; actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("expected no error, but got %q", err)
		}
	})

	t.Run("delete key from store", func(t *testing.T) {
		if actual, err := intr.Exec("DEL foo"); err == nil {
			if expected := 1; actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("expected no error, but got %q", err)
		}
	})

	t.Run("get deleted key from store", func(t *testing.T) {
		if actual, err := intr.Exec("GET foo"); err == nil {
			if actual != nil {
				t.Errorf("expected nil, got %v", actual)
			}
		} else {
			t.Errorf("expected no error, but got %q", err)
		}
	})

	t.Run("increment key from store", func(t *testing.T) {
		intr.Exec("INCR bar")
		if actual, err := intr.Exec("INCR bar"); err == nil {
			if expected := 2; actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("expected no error, but got %q", err)
		}
	})

	t.Run("set sorte set item to key", func(t *testing.T) {
		if actual, err := intr.Exec("ZADD fizz 3 three"); err == nil {
			if expected := 1; actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("expected no error, but got %q", err)
		}
	})

	t.Run("get sorted set cardinality", func(t *testing.T) {
		if actual, err := intr.Exec("ZCARD fizz"); err == nil {
			if expected := 1; actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("expected no error, but got %q", err)
		}
	})

	t.Run("get sorted set item rank", func(t *testing.T) {
		if actual, err := intr.Exec("ZRANK fizz three"); err == nil {
			if expected := 0; actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("expected no error, but got %q", err)
		}
	})

	t.Run("get sorted set range", func(t *testing.T) {
		if actual, err := intr.Exec("ZRANGE fizz 0 0"); err == nil {
			if expected := []string{"three"}; !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("expected no error, but got %q", err)
		}
	})

	t.Run("execute invalid command", func(t *testing.T) {
		if _, err := intr.Exec("SEY foo"); err == nil {
			t.Errorf("expected error, but got nil")
		}
	})
}
