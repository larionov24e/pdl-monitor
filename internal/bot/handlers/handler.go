package handlers

import (
	"context"
	"vdng-bot/internal/bot/session"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler interface {
	Handle(ctx context.Context, update tgbotapi.Update, session *session.UserSession) (tgbotapi.Chattable, error)
}
