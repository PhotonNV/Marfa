package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PhotonNV/Marfa/internal/application"
	"github.com/PhotonNV/Marfa/internal/infrastructure"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("MARFA_TOKEN")
	if token == "" {
		log.Fatal("MARFA_TOKEN is not set")
	}

	downloadDir := "marfa_downloads"
	if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create download directory: %v", err)
	}

	telegramBot, err := infrastructure.NewTelegramBot(token, nil, nil, nil)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}

	greetHandler := &application.GreetHandler{Bot: telegramBot}
	echoHandler := &application.EchoHandler{Bot: telegramBot}
	transmissionDownloader := infrastructure.NewTransmissionDownloader()

	telegramBot.GreetHandler = greetHandler
	telegramBot.EchoHandler = echoHandler
	telegramBot.TransmissionDownloader = transmissionDownloader

	updateChan := make(chan tgbotapi.Update)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go telegramBot.ListenForUpdates(ctx, updateChan)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Shutting down...")
		cancel()
	}()

	for update := range updateChan {
		if update.Message != nil {
			if update.Message.Text != "" {
				telegramBot.HandleTextMessage(update)
			} else {
				telegramBot.HandleFile(update)
			}
		}
	}
}
