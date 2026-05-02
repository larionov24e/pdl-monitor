package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vdng-bot/internal/bot"
	"vdng-bot/internal/scheduler"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	baseUrl := os.Getenv("BASE_URL")
	quit := make(chan os.Signal, 1)

	if token == "" {
		log.Panic("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	sch := scheduler.NewScheduler(60*time.Second, baseUrl)

	if sch == nil {
		slog.Error("Scheduler initialization failed")
		return
	}

	go sch.Run(ctx)

	api, err := bot.NewBot(token)

	if api == nil {
		slog.Error("Bot initialization failed")

		return
	}

	api.AddStorage(sch.ScheduleStorage)

	if err != nil {
		log.Panic(err)
	}

	go api.Start(ctx)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	}
}
