package helpers

import (
	"fmt"
	"time"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
)

type Summary struct {
	Participants []string
	MessageCount int64
	Messages     []Message
}

type Message struct {
	From string
	Time time.Time
	Text string
}

func ProcessSummary(messages []*models.Message) (*Summary, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages to process")
	}

	participantMap := make(map[string]bool)

	res := &Summary{
		MessageCount: int64(len(messages)),
		Messages:     make([]Message, 0, len(messages)),
	}

	for _, msg := range messages {
		if msg.Username != nil && *msg.Username != "" {
			participantMap[*msg.Username] = true
		}

		res.Messages = append(res.Messages, Message{
			From: safeString(msg.Username),
			Time: msg.CreatedAt,
			Text: safeString(msg.Text),
		})
	}

	res.Participants = make([]string, 0, len(participantMap))
	for username := range participantMap {
		res.Participants = append(res.Participants, username)
	}

	return res, nil
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
