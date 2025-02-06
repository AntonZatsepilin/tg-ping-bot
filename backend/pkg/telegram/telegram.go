package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var chatID int64

func Init(token string, chat int64) {
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panicf("Failed to create bot: %v", err)
	}

	chatID = chat
}

func SendMessage(msg string) {
	message := tgbotapi.NewMessage(chatID, msg)
	_, err := bot.Send(message)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	} else {
		log.Printf("Message sent: %s", msg)
	}
}

func HandleUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bot started. I'm monitoring the server.")
				bot.Send(msg)
			}
		}
	}
}
