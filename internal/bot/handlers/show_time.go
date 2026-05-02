package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"
	"vdng-bot/internal/bot/session"
	"vdng-bot/internal/scheduler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ShowCourts struct {
	storage *scheduler.Storage
}

func NewShowCourts(s *scheduler.Storage) *ShowCourts {
	return &ShowCourts{storage: s}
}

func (s *ShowCourts) Handle(_ context.Context, update tgbotapi.Update, _ *session.UserSession) (tgbotapi.Chattable, error) {
	_, payload, _ := strings.Cut(update.CallbackQuery.Data, ":")

	selected, err := time.Parse("02-01-2006", payload)
	if err != nil {
		return nil, fmt.Errorf("parse date %q: %w", payload, err)
	}

	s.storage.AddCache(update.CallbackQuery.Message.Chat.ID, scheduler.DateKey, selected.Format("2006-01-02"))

	var rows [][]tgbotapi.InlineKeyboardButton

	startTime, _ := time.Parse("15:04:05", "09:00:00")
	endTime, _ := time.Parse("15:04:05", "22:00:00")

	for t := startTime; !t.After(endTime); t = t.Add(time.Hour) {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				t.Format("15:04:05"),
				fmt.Sprintf("start_time_from:%s", t.Format("15:04")),
			),
		))
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please, select time from: ")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	return msg, nil
}
