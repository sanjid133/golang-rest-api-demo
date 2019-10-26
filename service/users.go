package service

import (
	"encoding/json"
	"github.com/sanjid133/rest-user-store/model"
	"net/http"
)

func (s *Service) PostUsers(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	var usr model.PostUser
	_ = json.NewDecoder(r.Body).Decode(usr)

	usrMdl := model.User{
		FirstName: usr.FirstName,
		LastName: usr.LastName,
		Password: usr.Password,
	}

	if err := s.usrRepo.Add(&usrMdl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(usrMdl)
}
