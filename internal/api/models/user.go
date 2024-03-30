package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"telegram_username"` //Telegram username
}
