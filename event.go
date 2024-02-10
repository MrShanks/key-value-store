package main

type Event struct {
	Sequence  int64
	EventType EventType
	Key       string
	Value     string
}

type EventType byte

const (
	EventDelete EventType = 1
	EventPut    EventType = 2
)
