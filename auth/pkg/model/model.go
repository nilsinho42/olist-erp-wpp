package model

// Model defined the structure of the token
type Token struct {
	ID         int    `json: "id"`
	Key        string `json: "key"`
	Lastupdate string `json: "lastupdate"`
}
