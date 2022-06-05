package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

type Todo struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	IsComplete bool   `json:"isComplete"`
}

type todoHandlers struct {
	sync.Mutex
	store map[string]Todo
}

func newTodoHandlers() *todoHandlers {
	return &todoHandlers{
		store: map[string]Todo{
			"id1": Todo{
				ID:         "id1",
				Title:      "Learn Golang",
				IsComplete: false,
			},
		},
	}
}

func (h *todoHandlers) restHandlers(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		h.get(writer, req)
		return
	case "POST":
		h.post(writer, req)
		return
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Method not allowed."))
		return
	}
}

func errorHandler(err error, writer http.ResponseWriter, req *http.Request) {
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
}

func (h *todoHandlers) post(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	errorHandler(err, writer, req)

	var todo Todo
	err = json.Unmarshal(bodyBytes, &todo)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	h.Lock()
	h.store[todo.ID] = todo
	defer h.Unlock()
}

func (h *todoHandlers) get(writer http.ResponseWriter, req *http.Request) {
	todos := make([]Todo, len(h.store))
	h.Lock()
	i := 0
	for _, todo := range h.store {
		todos[i] = todo
		i++
	}
	h.Unlock()
	jsonBytes, err := json.Marshal(todos)
	errorHandler(err, writer, req)
	writer.Header().Add("content-type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonBytes)
}

func main() {
	todoHandlers := newTodoHandlers()
	http.HandleFunc("/todos", todoHandlers.restHandlers)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
