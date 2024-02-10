package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var logger *FileTransactionLogger
var store = NewStore()

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi, I'm a key value store"))
}

func handlerGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := store.Get(key)
	if err != nil {
		log.Printf("key: '%s' %s", key, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("Value of '%s' is '%s'", key, value)

	w.WriteHeader(http.StatusFound)
	w.Write([]byte(value))
}

func handlerPut(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error parsing request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var kv KeyValue

	err = json.Unmarshal(body, &kv)
	if err != nil {
		log.Printf("Error Unmarshalling request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	store.Put(kv.Key, kv.Value)
	logger.WritePut(kv.Key, kv.Value)
	log.Printf("A new key value pair is added: %+v", kv)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("key: '%s' has been added with value: '%s'", kv.Key, kv.Value)))
}

func handlerDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	err := store.Delete(key)
	if err != nil {
		log.Printf("key: '%s' %s", key, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	logger.WriteDelete(key)

	log.Printf("key: %s has been deleted", key)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf("key: '%s' has been deleted", key)))
}

func InitFileTransactionLogger() {
	var err error
	logger, err = NewFileTransactionLogger("transaction.log")
	if err != nil {
		log.Fatalf("Error creating file: %s", err)
	}

	events := logger.ReadEvents()
	e, ok := Event{}, true

	for ok {
		e, ok = <-events
		switch e.EventType {
		case EventPut:
			store.Put(e.Key, e.Value)
		case EventDelete:
			store.Delete(e.Key)
		}

	}

	logger.Run()
}

func main() {
	InitFileTransactionLogger()

	r := mux.NewRouter()

	r.HandleFunc("/", handlerRoot)
	r.HandleFunc("/v1/key", handlerPut).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", handlerGet).Methods("GET")
	r.HandleFunc("/v1/key/{key}", handlerDelete).Methods("DELETE")

	log.Printf("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
