package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func (b *Bot) Boot(botToken string) {
	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	b.bot = bot
	go b.Watch()
}

func (b *Bot) Watch() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			b.bot.Send(msg)
		}
	}
}

func (b *Bot) SendPhotoAlarm(bytesData []byte) bool {
	fmt.Println(len(bytesData))

	photoFileBytes := tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: bytesData,
	}

	photoMessage := tgbotapi.NewPhoto(int64(-1001317502508), photoFileBytes)
	photoMessage.File = photoFileBytes

	_, err := b.bot.Send(photoMessage)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func (b *Bot) SendVideoAlarm() bool {
	return true
}

func getNumericKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔇", "muteCommand"),
			tgbotapi.NewInlineKeyboardButtonData("🔉", "unMuteCommand"),
		),
	)
}
