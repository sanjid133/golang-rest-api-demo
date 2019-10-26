package service

import "github.com/gorilla/mux"

func (s *Service) Route() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users", s.PostUsers).Methods("POST")
	r.HandleFunc("/user/{id}", s.GetUser).Methods("GET")

	return r
}
