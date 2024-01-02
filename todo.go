package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	IsDone      bool               `json:"is_done"`
}

type Response struct {
	Body interface{} `json:"body,omitempty"`
	Err  error       `json:"error,omitempty"`
}

func NewResponse(body any, err error) *Response {
	return &Response{
		Body: body,
		Err:  err,
	}
}

// TODO: Add logger
// TODO: Handle errors better
// TODO: Refactor code and make it leaner

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

func (s *Server) GetTodoFromQuery(w http.ResponseWriter, r *http.Request) {
	coll := s.client.Database("lazydev").Collection("tasks")

	var result Task

	id := r.URL.Query().Get("id")

	fmt.Println(id)

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewResponse(nil, errors.New("invalid id")))
		return
	}

	oid, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{{Key: "_id", Value: oid}}
	err := coll.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewResponse(nil, errors.New("error while fetching data")))
	}

	json.NewEncoder(w).Encode(NewResponse(result, nil))
}

func (s *Server) GetTodoFromPathVariable(w http.ResponseWriter, r *http.Request) {
	coll := s.client.Database("lazydev").Collection("tasks")

	var result Task

	id := mux.Vars(r)["id"] // for path variable

	fmt.Println(id)

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewResponse(nil, errors.New("invalid id")))
		return
	}

	oid, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{{Key: "_id", Value: oid}}
	err := coll.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewResponse(nil, errors.New("error while fetching data")))
	}

	json.NewEncoder(w).Encode(NewResponse(result, nil))
}

func (s *Server) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	coll := s.client.Database("lazydev").Collection("tasks")

	var tasks []Task

	filter := bson.D{{}}
	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewResponse(nil, err))
		return
	}

	if err = cursor.All(context.Background(), &tasks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewResponse(nil, err))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(tasks, nil))

}

func (s *Server) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	coll := s.client.Database("lazydev").Collection("tasks")

	var newTask Task

	json.NewDecoder(r.Body).Decode(&newTask)

	filter := bson.D{{Key: "_id", Value: newTask.Id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isdone", Value: newTask.IsDone}}}}
	_, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewResponse(nil, err))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(newTask, nil))
}
