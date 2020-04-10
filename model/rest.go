package model

type LoginRequest struct {
	ClientID string `json:"client_id"`
	UserName string `json:"user_name"`
	Color    string `json:"color"`
}
