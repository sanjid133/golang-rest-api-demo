package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sanjid133/rest-user-store/model"
	"github.com/sanjid133/rest-user-store/util"
	"net/http"
	"strings"
	"time"
)

func (s *Service) PostTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tags model.PostTags
	err := json.NewDecoder(r.Body).Decode(&tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	if _, found := params["id"]; !found {
		http.Error(w, errors.Errorf("id is requried").Error(), http.StatusBadRequest)
		return
	}
	userID := params["id"]

	usr, err := s.usrRepo.Get(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if usr == nil {
		http.Error(w, errors.Errorf("User not found with id %v", userID).Error(), http.StatusBadRequest)
		return
	}

	for _, t := range tags.Tags {
		tagModel := model.Tag{
			Tag:      t,
			UserID:   usr.ID,
			ExpireAt: util.Now().Add(time.Duration(tags.Expiry) * time.Millisecond),
		}
		if err := s.tagRepo.Add(&tagModel); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func (s *Service) SearchByTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if _, found := params["tags"]; !found {
		http.Error(w, errors.Errorf("id is requried").Error(), http.StatusBadRequest)
		return
	}
	sTags := strings.Split(params["tags"], ",")

	tags, err := s.tagRepo.ListTags(sTags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userIDs := []string{}

	for _, t := range tags {
		userIDs = append(userIDs, t.UserID)
	}
	users, err := s.usrRepo.ListUsers(userIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := []RespTagUser{}
	for _, u := range users {
		tags, err := s.tagRepo.ListByUserID(u.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tu := RespTagUser{
			ID:   u.ID,
			Name: fullName(u.FirstName, u.LastName),
			Tags: []string{},
		}
		for _, t := range tags {
			tu.Tags = append(tu.Tags, t.Tag)
		}
		resp = append(resp, tu)
	}
	json.NewEncoder(w).Encode(RespUsers{resp})

}
