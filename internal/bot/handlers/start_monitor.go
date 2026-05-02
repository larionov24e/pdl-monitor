package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"time"
	"vdng-bot/internal/bot/session"
	"vdng-bot/internal/court"
	scheduler2 "vdng-bot/internal/scheduler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartMonitor struct {
	storage *scheduler2.Storage
	api     *tgbotapi.BotAPI
}

func NewStartMonitor(storage *scheduler2.Storage, api *tgbotapi.BotAPI) *StartMonitor {
	return &StartMonitor{
		storage: storage,
		api:     api,
	}
}

func (sm *StartMonitor) Handle(ctx context.Context, update tgbotapi.Update, sess *session.UserSession) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Monitor is started. We will notice you for new available courts soon")

	if sess.CancelFunc != nil {
		sess.CancelFunc()
	}
	monitorCtx, cancel := context.WithCancel(ctx)
	sess.CancelFunc = cancel

	go func(ctx context.Context, s *session.UserSession, storage *scheduler2.Storage, api *tgbotapi.BotAPI) {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

	mainStep:
		for {
			currentDate := time.Now().Format("2006-01-02")

			select {
			case <-monitorCtx.Done():
				break mainStep
			case <-ctx.Done():
				break mainStep
			case <-ticker.C:
				var resultAvailableCourts []court.AvailableCourt
				var actualCourts []scheduler2.Court

				allCourts := storage.GetCourtsByDate(currentDate)

				if len(allCourts) == 0 {
					slog.Error("GetCourts not found for ", "date", currentDate)
					continue
				}

				availableCourtsFromSession := storage.CacheByKey(s.SessionId, scheduler2.SnapshotKey)

				err := json.Unmarshal([]byte(availableCourtsFromSession), &resultAvailableCourts)
				if err != nil {
					slog.Error("Unmarshal err", slog.Any("err", err))
					continue
				}

				newCourtsSnapshot := court.FilterByTimeRange(s.TimeFrom, s.TimeTo, actualCourts)

				var foundCourts []court.AvailableCourt

				for _, newCourt := range newCourtsSnapshot {
					for _, oldCourt := range resultAvailableCourts {
						if newCourt.Name != oldCourt.Name {
							continue
						}

						var newSlots []string
						for _, slot := range newCourt.Time {
							if !slices.Contains(oldCourt.Time, slot) {
								newSlots = append(newSlots, slot)
							}
						}

						if len(newSlots) > 0 {
							foundCourts = append(foundCourts, court.AvailableCourt{
								Name: newCourt.Name,
								Time: newSlots,
							})
						}
					}
				}

				if len(foundCourts) > 0 {
					text := "New courts available:\n"
					for _, c := range foundCourts {
						text += fmt.Sprintf("- %s: %v\n", c.Name, c.Time)
					}

					notify := tgbotapi.NewMessage(s.SessionId, text)

					rows := make([][]tgbotapi.InlineKeyboardButton, 0, 1)

					rows = append(rows, tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(
							"Continue",
							"continue_monitor:",
						),
					))

					notify.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

					_, err := api.Send(notify)
					if err != nil {
						return
					}

					select {
					case <-monitorCtx.Done():
						break mainStep
					case <-ctx.Done():
						return
					case <-s.WaitContinue():
						ticker.Reset(5 * time.Minute)
						continue
					}
				}
			}
		}

	}(monitorCtx, sess, sm.storage, sm.api)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 1)

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"Stop monitor",
			"stop_monitor:",
		),
	))

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}
