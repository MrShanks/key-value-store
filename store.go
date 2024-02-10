package main

import "errors"

var ErrorKeyNotFound = errors.New("not found")

type Store struct {
	m map[string]string
}

type KeyValue struct {
	Key   string
	Value string
}

func NewStore() *Store {
	return &Store{
		m: make(map[string]string),
	}
}

func (s *Store) Get(key string) (string, error) {
	value, ok := s.m[key]
	if !ok {
		return "", ErrorKeyNotFound
	}
	return value, nil
}

func (s *Store) Put(key, value string) {
	s.m[key] = value
}

func (s *Store) Delete(key string) error {
	_, ok := s.m[key]
	if !ok {
		return ErrorKeyNotFound
	}
	delete(s.m, key)
	return nil
}
