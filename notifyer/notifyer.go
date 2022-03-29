package notifyer

import (
	"github.com/alexdin/tinygonvr/notifyer/telegram"
)

type NotifyTransfer interface {
	Boot(string)
	SendPhotoAlarm([]byte) bool
	SendVideoAlarm() bool
}

type Notify struct {
	BotToken  string
	ChannelId int
	BotType   int
	PhotoChan chan []byte
}

var botInstance NotifyTransfer = nil
var config Notify

// todo implements more bots for notify
const ChannelTelegram = 0

func Boot(notify Notify) {
	config = notify
	switch config.BotType {
	case ChannelTelegram:
		botInstance = new(telegram.Bot)
	}

	botInstance.Boot(config.BotToken)
	go waitForPhoto()
}

func waitForPhoto() {
	botInstance.SendPhotoAlarm(<-config.PhotoChan)
}

func PutPhotoToChannel(data []byte) {
	config.PhotoChan <- data
}
