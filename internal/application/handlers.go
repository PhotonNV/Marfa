package application

import (
	"fmt"
	"log"

	"github.com/PhotonNV/Marfa/internal/infrastructure"
)

type GreetHandler struct {
	Bot *infrastructure.TelegramBot // Ссылка на TelegramBot
}

func (h *GreetHandler) Handle(cmd GreetCommand) {
	response := "Добро пожаловать в Marfa!"
	log.Printf("Sending greeting to chat %d: %s", cmd.ChatID, response)
	h.Bot.SendMessage(cmd.ChatID, response) // Отправка сообщения в чат
}

type EchoHandler struct {
	Bot *infrastructure.TelegramBot // Ссылка на TelegramBot
}

func (h *EchoHandler) Handle(cmd EchoCommand) {
	response := fmt.Sprintf("Вы написали: %s", cmd.Message)
	log.Printf("Sending echo to chat %d: %s", cmd.ChatID, response)
	h.Bot.SendMessage(cmd.ChatID, response) // Отправка сообщения в чат
}
