package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/config"
)

type Bot struct {
	conn   string
	client *http.Client
}

func NewBot(cfg config.Telegram) *Bot {
	return &Bot{
		conn:   fmt.Sprintf("api.telegram.org/bot%s/", cfg.Token),
		client: &http.Client{},
	}
}

type Message struct {
	Username string `json:"chat_id"`
	Text     string `json:"text"`
}

func (tg *Bot) SendNotification(username, message string) error {
	msg := Message{Username: username, Text: message}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return fmt.Errorf("failed to encode notification: %w", err)
	}

	if _, err := tg.client.Post(tg.conn+"sendMessage", "application/json", &buf); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}
