package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Service struct {
	botToken string
	chatID   string
	client   *http.Client
}

func NewService(botToken, chatID string) *Service {
	return &Service{
		botToken: botToken,
		chatID:   chatID,
		client:   http.DefaultClient,
	}
}

type message struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func (s *Service) SendNotification(text string) error {
	msg := message{
		ChatID: s.chatID,
		Text:   text,
	}

	jsonBody, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("telegram API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
