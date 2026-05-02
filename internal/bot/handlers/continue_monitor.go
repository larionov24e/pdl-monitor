package handlers

import (
	"context"
	"vdng-bot/internal/bot/session"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ContinueMonitor struct{}

func NewContinueMonitor() *ContinueMonitor {
	return &ContinueMonitor{}
}

func (cm *ContinueMonitor) Handle(_ context.Context, _ tgbotapi.Update, s *session.UserSession) (tgbotapi.Chattable, error) {
	s.Continue()

	msg := tgbotapi.NewMessage(s.SessionId, "Monitoring continued. We will notify you about new courts.")
	return msg, nil
}
