package model

type JWTToken struct {
	Token    string `json:"token"`
	Describe string `json:"describe"`
}
