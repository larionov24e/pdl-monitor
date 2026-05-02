package bot

import (
	"fmt"
	"strings"
	"vdng-bot/internal/bot/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	commands  map[string]handlers.Handler
	callbacks map[string]handlers.Handler
}

func NewRouter() *Router {
	return &Router{
		commands:  make(map[string]handlers.Handler),
		callbacks: make(map[string]handlers.Handler),
	}
}

func (r *Router) Command(cmd string, h handlers.Handler) {
	r.commands[cmd] = h
}

func (r *Router) Callback(prefix string, h handlers.Handler) {
	r.callbacks[prefix] = h
}

func (r *Router) Resolve(update tgbotapi.Update) handlers.Handler {
	switch {
	case update.Message != nil:
		msgChecker := checkMessageType(&update)

		if msgChecker != "" {
			return r.commands[msgChecker]
		}

		return r.commands[update.Message.Text]
	case update.CallbackQuery != nil:
		prefix, _, _ := strings.Cut(update.CallbackQuery.Data, ":")
		fmt.Println(prefix)
		return r.callbacks[prefix]
	}
	return nil
}

func checkMessageType(update *tgbotapi.Update) string {
	if update.Message.ReplyToMessage != nil {
		return "start_time_from"
	}

	return ""
}
