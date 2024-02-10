package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewStore(t *testing.T) {
	t.Run("Create empty Store", func(t *testing.T) {
		got := NewStore()
		expected := make(map[string]string)

		if !reflect.DeepEqual(got.m, expected) {
			t.Errorf("Expected %s map got: %s", expected, got.m)
		}
	})
}

func TestPut(t *testing.T) {
	store := NewStore()
	t.Run("When a new Value is store should be able to get it with its key", func(t *testing.T) {
		store.Put("key", "value")
		got := store.m["key"]
		expected := "value"

		if got != expected {
			t.Errorf("Got: %s expected: %s", got, expected)
		}
	})
}

func TestGet(t *testing.T) {
	store := NewStore()
	t.Run("Get should return error not found when retrieving a key that doesn't exist", func(t *testing.T) {
		_, err := store.Get("key")
		expected := errors.New("not found")

		if expected.Error() != err.Error() {
			t.Errorf("expected: %s, got: %s", expected, err)
		}
	})

	t.Run("Get should return value when key exists", func(t *testing.T) {
		store.m["key"] = "value"
		got, _ := store.Get("key")
		expected := "value"

		if expected != got {
			t.Errorf("expected: %s, got: %s", expected, got)
		}
	})
}

func TestDelete(t *testing.T) {
	store := NewStore()
	t.Run("Get should return error not found when retrieving a key that doesn't exist", func(t *testing.T) {
		err := store.Delete("key")
		expected := errors.New("not found")

		if expected.Error() != err.Error() {
			t.Errorf("expected: %s, got: %s", expected, err)
		}
	})

	t.Run("When deleting a key value pair it cannot be retrieved anymore", func(t *testing.T) {
		store.m["key"] = "value"
		store.Delete("key")

		got := store.m["key"]
		expected := ""

		if got != expected {
			t.Errorf("Got: %s expected: %s", got, expected)
		}
	})
}
