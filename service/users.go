package service

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sanjid133/rest-user-store/model"
	"golang.org/x/crypto/pbkdf2"
	"net/http"
)

func (s *Service) PostUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var usr model.PostUser
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usrMdl := &model.User{
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
	}
	if err := SetPassword(usrMdl, usr.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.usrRepo.Add(usrMdl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := respPostUsers{
		ID: usrMdl.ID,
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *Service) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if _, found := params["id"]; !found {
		http.Error(w, errors.Errorf("id is requried").Error(), http.StatusBadRequest)
		return
	}
	usr, err := s.usrRepo.Get(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if usr == nil {
		http.Error(w, errors.Errorf("User not found with id %v", params["id"]).Error(), http.StatusBadRequest)
		return
	}
	resp := respGetUser{
		ID:   usr.ID,
		Name: fmt.Sprintf("%s %s", usr.FirstName, usr.LastName),
	}
	json.NewEncoder(w).Encode(resp)
}

func SetPassword(u *model.User, pass string) error {
	u.PassSalt = make([]byte, 16)
	if _, err := rand.Read(u.PassSalt); err != nil {
		return err
	}
	u.PassAlgo = model.AlgoSha256
	u.PassIter = 10000
	u.PassHash = pbkdf2.Key([]byte(pass), u.PassSalt, u.PassIter, 32, sha256.New)
	return nil
}

// VerifyUserPassword verifies if pass matcher to user u
func VerifyUserPassword(u *model.User, pass string) (bool, error) {
	hash := pbkdf2.Key([]byte(pass), u.PassSalt, u.PassIter, 32, sha256.New)
	eq := bytes.Equal(u.PassHash, hash)
	if !eq {
		return false, errors.Errorf("Invalid login cred")
	}
	return true, nil

}
