package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type Task struct {
	Title       string
	Description string
	IsDone      bool
}

type Response struct {
	Body interface{}
	Err  error
}

func NewResponse(body any, err error) Response {
	return Response{
		Body: body,
		Err:  err,
	}
}

func (s *Server) CreateTodo(w http.ResponseWriter, r *http.Request) {
	coll := s.client.Database("lazydev").Collection("tasks")

	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)

	res, err := coll.InsertOne(context.Background(), task)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewResponse(nil, err))
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(NewResponse(res, nil))
}
