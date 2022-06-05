package main

import (
	"encoding/json"
	"net/http"
)

type Todo struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	IsComplete bool   `json:"isComplete"`
}

type todoHandlers struct {
	store map[string]Todo
}

func newTodoHandlers() *todoHandlers {
	return &todoHandlers{
		store: map[string]Todo{
			"id1": {
				ID:         "id1",
				Title:      "Learn Golang",
				IsComplete: false,
			},
		},
	}
}

func (h *todoHandlers) getAll(writer http.ResponseWriter, req *http.Request) {
	todos := make([]Todo, len(h.store))
	i := 0
	for _, todo := range h.store {
		todos[i] = todo
		i++
	}
	jsonBytes, err := json.Marshal(todos)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		panic(err)
	}
	writer.Header().Add("content-type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonBytes)
}

func main() {
	todoHandlers := newTodoHandlers()
	http.HandleFunc("/todos", todoHandlers.getAll)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
