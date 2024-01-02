package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	client *mongo.Client
}

func NewServer(c *mongo.Client) *Server {
	return &Server{
		client: c,
	}
}

func (s *Server) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/todo", s.CreateTodo).Methods(http.MethodPost)
	r.HandleFunc("/todo", s.GetTodoFromQuery).Methods(http.MethodGet)
	r.HandleFunc("/todo/{id}", s.GetTodoFromPathVariable).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    ":7000",
		Handler: r,
	}

	fmt.Println("running on port 7000...")
	log.Fatal(srv.ListenAndServe())
}
