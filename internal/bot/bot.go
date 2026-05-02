package bot

import (
	"context"
	"log"
	"log/slog"
	"sync"
	handlers2 "vdng-bot/internal/bot/handlers"
	"vdng-bot/internal/bot/session"
	"vdng-bot/internal/scheduler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	sessions map[int64]*session.UserSession
	mutex    sync.Mutex
	storage  *scheduler.Storage
}

func (b *Bot) GetOrCreateSession(id int64) *session.UserSession {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.sessions) == 100 {
		slog.Error("GetOrCreateSession called more than 100 sessions")
		return nil
	}

	if s, ok := b.sessions[id]; ok {
		return s
	}
	s := session.NewUserSession(id)
	b.sessions[id] = s
	return s
}

func (bot *Bot) AddStorage(sch *scheduler.Storage) {
	bot.mutex.Lock()
	defer bot.mutex.Unlock()
	bot.storage = sch
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		slog.Error("error msg ", "err", err.Error())
		return nil, err
	}

	api.Debug = true

	log.Printf("Authorized on account %s", api.Self.UserName)
	return &Bot{api: api, sessions: make(map[int64]*session.UserSession)}, nil
}

func (b *Bot) Start(ctx context.Context) {
	router := NewRouter()
	router.Command("/start", handlers2.NewStart())
	router.Callback("start_time_from", handlers2.NewTimeTo(b.storage))
	router.Callback("start_time_to", handlers2.NewSelectTime(b.storage))
	router.Callback("stop_monitor", handlers2.NewStopMonitor())
	router.Callback("date", handlers2.NewSelectDate())
	router.Callback("selected_date", handlers2.NewShowCourts(b.storage))
	router.Callback("start_monitor", handlers2.NewStartMonitor(b.storage, b.api))
	router.Callback("continue_monitor", handlers2.NewContinueMonitor())

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	defer updates.Clear()
	var sessionId int64

	for update := range updates {
		go func(update tgbotapi.Update, bot *Bot, sessionId int64) {
			h := router.Resolve(update)
			if h == nil {
				return
			}

			if update.Message != nil {
				sessionId = update.Message.Chat.ID
			} else if update.CallbackQuery != nil {
				sessionId = update.CallbackQuery.Message.Chat.ID
			}

			sess := b.GetOrCreateSession(sessionId)

			msg, err := h.Handle(ctx, update, sess)
			if err != nil {
				slog.Error("handler error", slog.Any("err", err))
				return
			}
			if msg == nil {
				return
			}
			if _, err := b.api.Send(msg); err != nil {
				slog.Error("send error", slog.Any("err", err))
			}
		}(update, b, sessionId)
	}
}
