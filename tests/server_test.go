package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sanjid133/rest-user-store/model"
	"github.com/sanjid133/rest-user-store/service"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

const serverURL = "http://localhost:8000"

func TestServer(t *testing.T) {
	userID := Service_PostUsers()
	Service_GetUser(userID)
	Service_PostTags(userID)
	Service_GetUserByTag([]string{"tag1", "tag2"})
}

func Service_GetUserByTag(tags []string) {
	req, _ := http.NewRequest("GET", serverURL+"/users?tags="+strings.Join(tags, ","), nil)
	req.Header.Set("Content-Type", "application/json")

	resp := executeRequest(req)
	fmt.Println("Get user response", resp)

	ru := service.RespUsers{}
	if err := json.Unmarshal([]byte(resp), &ru); err != nil {
		panic(err)
	}
	fmt.Println(ru)

}

func Service_PostTags(uid string) {
	data := model.PostTags{
		Tags:   []string{"tag1", "tag2"},
		Expiry: 9000,
	}
	pd, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(pd))
	req, _ := http.NewRequest("POST", serverURL+"/user/"+uid+"/tags", bytes.NewBuffer(pd))
	req.Header.Set("Content-Type", "application/json")

	resp := executeRequest(req)
	fmt.Println("Get user tags response", resp)

}

func Service_GetUser(uid string) {
	req, _ := http.NewRequest("GET", serverURL+"/user/"+uid, nil)
	req.Header.Set("Content-Type", "application/json")

	resp := executeRequest(req)
	fmt.Println("Get user response", resp)

	ru := service.RespGetUser{}
	if err := json.Unmarshal([]byte(resp), &ru); err != nil {
		panic(err)
	}

}

func Service_PostUsers() string {
	data := model.PostUser{
		FirstName: "Sanjidul",
		LastName:  "Hoque",
		Password:  "password",
	}
	pd, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("POST", serverURL+"/users", bytes.NewBuffer(pd))
	req.Header.Set("Content-Type", "application/json")

	resp := executeRequest(req)
	fmt.Println("Get user response", resp)

	ru := service.RespPostUsers{}
	if err := json.Unmarshal([]byte(resp), &ru); err != nil {
		panic(err)
	}
	return ru.ID

}

func executeRequest(req *http.Request) string {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}
