package handlers

import (
	"context"
	"vdng-bot/internal/bot/session"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StopMonitor struct {
}

func NewStopMonitor() *StopMonitor {
	return &StopMonitor{}
}

func (sm *StopMonitor) Handle(ctx context.Context, update tgbotapi.Update, s *session.UserSession) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(s.SessionId, "Session doesnt exist")

	if s.CancelFunc != nil {
		s.CancelFunc()
		msg = tgbotapi.NewMessage(s.SessionId, "Done. Write /start for new session")
	}

	return msg, nil
}
