package service

type RespPostUsers struct {
	ID string `json:"id"`
}

type RespPostTag struct {}

type RespGetUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RespTagUser struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type RespUsers struct {
	Users []RespTagUser `json:"users"`
}
