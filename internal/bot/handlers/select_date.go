package handlers

import (
	"context"
	"log/slog"
	"time"
	"vdng-bot/internal/bot/session"
	"vdng-bot/internal/scheduler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SelectDate struct{}

func NewSelectDate() *SelectDate { return &SelectDate{} }

func (s *SelectDate) Handle(_ context.Context, update tgbotapi.Update, _ *session.UserSession) (tgbotapi.Chattable, error) {
	dates := scheduler.GenerateDates()
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(dates))

	for _, date := range dates {
		d, err := time.Parse("2006-01-02", date)
		if err != nil {
			slog.Warn("failed to parse date", slog.String("date", date))
			continue
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				d.Format("Monday 02/01/2006"),
				"selected_date:"+d.Format("02-01-2006"),
			),
		))
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Select Option")
	if len(rows) > 0 {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	}
	return msg, nil
}
