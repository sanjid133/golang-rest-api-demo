package service

type respPostUsers struct {
	ID string `json:"id"`
}

type respGetUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
