package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type FileTransactionLogger struct {
	Event        chan Event
	LastSequence int64
	File         *os.File
}

func NewFileTransactionLogger(filename string) (*FileTransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log: %s", err)
	}

	return &FileTransactionLogger{File: file, Event: make(chan Event, 16)}, nil
}

func (f *FileTransactionLogger) WritePut(key, value string) {
	f.Event <- Event{
		EventType: EventPut,
		Key:       key,
		Value:     value,
	}
}

func (f *FileTransactionLogger) WriteDelete(key string) {
	f.Event <- Event{
		EventType: EventDelete,
		Key:       key,
	}
}

func (f *FileTransactionLogger) Run() {
	go func() {
		for e := range f.Event {
			f.LastSequence++
			fmt.Fprintf(f.File, "%d\t%d\t%s\t%s\n", f.LastSequence, e.EventType, e.Key, e.Value)
		}
	}()
}

func (f *FileTransactionLogger) ReadEvents() chan Event {
	log.Printf("Reading past events from file...")
	scanner := bufio.NewScanner(f.File)
	outEvent := make(chan Event, 16)

	go func() {
		defer close(outEvent)
		for scanner.Scan() {
			var e Event
			line := scanner.Text()
			words := strings.Fields(line)
			if len(words) == 3 {
				fmt.Sscanf(
					line,
					"%d\t%d\t%s\t\n", &e.Sequence, &e.EventType, &e.Key)

			} else {
				fmt.Sscanf(
					line,
					"%d\t%d\t%s\t%s\n", &e.Sequence, &e.EventType, &e.Key, &e.Value)
			}

			log.Printf("Loading event %+v", e)
			outEvent <- e
		}
	}()
	return outEvent
}
