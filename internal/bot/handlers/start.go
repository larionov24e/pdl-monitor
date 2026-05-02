package handlers

import (
	"context"
	"vdng-bot/internal/bot/session"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Start struct{}

func NewStart() *Start { return &Start{} }

func (s *Start) Handle(ctx context.Context, update tgbotapi.Update, session *session.UserSession) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Select Option")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Please, select date", "date:select_date"),
		),
	)
	return msg, nil
}
