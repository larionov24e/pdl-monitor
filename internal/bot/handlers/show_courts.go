package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"vdng-bot/internal/bot/session"
	"vdng-bot/internal/court"
	"vdng-bot/internal/scheduler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SelectTime struct {
	storage *scheduler.Storage
}

func NewSelectTime(storage *scheduler.Storage) *SelectTime {
	return &SelectTime{
		storage: storage,
	}
}

func (t *SelectTime) Handle(_ context.Context, update tgbotapi.Update, s *session.UserSession) (tgbotapi.Chattable, error) {
	_, afterTime, ok := strings.Cut(update.CallbackQuery.Data, ":")

	if !ok {
		slog.Error("invalid callback query")
		return nil, fmt.Errorf(`invalid callback query`)
	}

	t.storage.AddCache(s.SessionId, scheduler.TimeToKey, afterTime)

	beforeTime := t.storage.CacheByKey(s.SessionId, scheduler.TimeFromKey)

	fmt.Println(afterTime, beforeTime)
	beforeTimeFormat, err := time.Parse("15:04", beforeTime)

	if err != nil {
		slog.Error(err.Error())
		return nil, fmt.Errorf(`invalid time format`)
	}

	afterTimeFormat, err := time.Parse("15:04", afterTime)

	if err != nil {
		slog.Error(err.Error())
	}

	s.TimeFrom = beforeTimeFormat
	s.TimeTo = afterTimeFormat

	msg := tgbotapi.NewMessage(s.SessionId, "Select timeframe")

	dateFromSessionStorage := t.storage.GetSessionCacheByKey(s.SessionId, scheduler.DateKey)
	courtsForCurrentDate := t.storage.GetCourtsByDate(dateFromSessionStorage)

	availableCourts := court.FilterByTimeRange(beforeTimeFormat, afterTimeFormat, courtsForCurrentDate)

	if len(availableCourts) > 0 {
		jsonAvailableCourts, err := json.Marshal(availableCourts)

		if err != nil {
			slog.Error(err.Error())
		}

		t.storage.AddCache(s.SessionId, scheduler.SnapshotKey, string(jsonAvailableCourts))
	}

	messageText := ""

	for _, course := range availableCourts {
		messageText += fmt.Sprintf("Name: %s \n", course.Name)
		for _, time := range course.Time {
			messageText += fmt.Sprintf("%s \n", time)
		}
	}

	msg.Text = messageText
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 1)

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"Start monitor",
			"start_monitor:",
		),
	))

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}
