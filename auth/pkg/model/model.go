package model

// Model defined the structure of the token
type Token struct {
	ID           int    `json:"id"`
	Key          string `json:"key"`
	Lastupdate   string `json:"lastupdate"`
	RefreshToken string `json:"refresh_token"`
}

type contextKey string

const ContextKey contextKey = "key"
const RefreshTokenKey contextKey = "refresh_token"

type RefreshResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}
