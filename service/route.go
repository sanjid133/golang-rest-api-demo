package service

import "github.com/gorilla/mux"

func (s *Service) Route() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users", s.PostUsers).Methods("POST")
	r.HandleFunc("/user/{id}", s.GetUser).Methods("GET")

	r.HandleFunc("/user/{id}/tags", s.PostTags).Methods("POST")
	r.HandleFunc("/users", s.SearchByTags).Queries("tags", "{tags}").Methods("GET")

	return r
}
