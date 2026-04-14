package models

type UserSession struct {
	UserID   int    `json:"id"`
	Username string `json:"username"`

	SessionID       string `json:"session_id"`
	IsAuthenticated bool   `json:"is_authenticated"`
}
