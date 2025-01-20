package infrastructure

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PhotonNV/Marfa/internal/application"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot                    *tgbotapi.BotAPI
	GreetHandler           *application.GreetHandler
	EchoHandler            *application.EchoHandler
	TransmissionDownloader *TransmissionDownloader
}

func NewTelegramBot(token string, greetHandler *application.GreetHandler, echoHandler *application.EchoHandler, downloader *TransmissionDownloader) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{bot: bot, GreetHandler: greetHandler, EchoHandler: echoHandler, TransmissionDownloader: downloader}, nil
}

func (tb *TelegramBot) HandleTextMessage(update tgbotapi.Update) {
	if update.Message.Text != "" {
		log.Printf("Received text message: %s", update.Message.Text)

		if update.Message.Text == "/start" {
			cmd := application.GreetCommand{ChatID: update.Message.Chat.ID}
			tb.GreetHandler.Handle(cmd)
		} else {
			cmd := application.EchoCommand{ChatID: update.Message.Chat.ID, Message: update.Message.Text}
			tb.EchoHandler.Handle(cmd)
		}
	}
}

func (tb *TelegramBot) HandleFile(update tgbotapi.Update) {
	if update.Message.Document != nil {
		fileID := update.Message.Document.FileID
		file, err := tb.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			log.Printf("Failed to get file: %v", err)
			return
		}
		log.Printf("Received file: %s", file.FilePath)

		filePath := filepath.Join("downloads", update.Message.Document.FileName) // Папка для сохранения
		if err := tb.downloadFile(file.FilePath, filePath); err != nil {
			log.Printf("Failed to download file: %v", err)
			return
		}

		if err := tb.TransmissionDownloader.Download(filePath); err != nil {
			log.Printf("Failed to download torrent: %v", err)
			return
		}

		log.Println("Torrent download started successfully.")
	}
}

func (tb *TelegramBot) downloadFile(filePath, destPath string) error {
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", tb.bot.Token, filePath)

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	return err
}

func (tb *TelegramBot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := tb.bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (tb *TelegramBot) ListenForUpdates(ctx context.Context, updateChan chan<- tgbotapi.Update) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tb.bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				if update.Message.Text != "" {
					tb.HandleTextMessage(update)
				} else {
					tb.HandleFile(update)
				}
				updateChan <- update
			}
		case <-ctx.Done():
			return
		}
	}
}
