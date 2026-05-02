package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"vdng-bot/internal/bot/session"
	"vdng-bot/internal/scheduler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TimeTo struct {
	storage *scheduler.Storage
}

func NewTimeTo(storage *scheduler.Storage) *TimeTo {
	return &TimeTo{
		storage: storage,
	}
}

func (tt *TimeTo) Handle(_ context.Context, update tgbotapi.Update, s *session.UserSession) (tgbotapi.Chattable, error) {
	_, timeFrom, ok := strings.Cut(update.CallbackQuery.Data, ":")

	if !ok {
		slog.Error("time_to: invalid callback query")
		return nil, fmt.Errorf(`time_to: invalid callback query`)
	}

	tt.storage.AddCache(s.SessionId, scheduler.TimeFromKey, timeFrom)

	var rows [][]tgbotapi.InlineKeyboardButton

	startTime, _ := time.Parse("15:04:05", "09:00:00")
	endTime, _ := time.Parse("15:04:05", "22:00:00")

	for t := startTime; !t.After(endTime); t = t.Add(time.Hour) {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				t.Format("15:04:05"),
				fmt.Sprintf("start_time_to:%s", t.Format("15:04")),
			),
		))
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please, select time to: ")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	return msg, nil
}
