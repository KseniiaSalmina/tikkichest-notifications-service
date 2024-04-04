package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/config"
	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/notifier"
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

func (tg *Bot) SendNotification(username string, event notifier.Event) error {
	msg := tg.messageFromEvent(username, event)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return fmt.Errorf("failed to encode notification: %w", err)
	}

	if _, err := tg.client.Post(tg.conn+"sendMessage", "application/json", &buf); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}

func (tg *Bot) messageFromEvent(username string, event notifier.Event) Message {
	text := fmt.Sprintf("Your %s №%d has been %s", string(event.Object), event.ObjectID, string(event.Change))
	return Message{
		Username: username,
		Text:     text}
}
