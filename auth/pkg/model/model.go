package model

// Model defined the structure of the token
type Token struct {
	id         int    `json: "id"`
	key        string `json: "key"`
	lastupdate string `json: "lastupdate"`
}
