package main

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
// #cgo CFLAGS: -std=c11 -g
// #include "../ffmpeg/ffmpeg.h"
import "C"
import (
	"fmt"
	"github.com/alexdin/tinygonvr/alarm"
	"github.com/alexdin/tinygonvr/ffmpeg"
	"github.com/alexdin/tinygonvr/notifyer"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {

	// load config data
	config := loadConfig()

	// boot config alert monitoring
	go alarm.Boot(config.getAlertConfig())
	go notifyer.Boot(
		notifyer.Notify{
			BotType:   notifyer.ChannelTelegram,
			BotToken:  config.Notify.BotToken,
			ChannelId: config.Notify.ChannelId,
			PhotoChan: make(chan []byte),
		},
	)

	ffmpeg.Boot(config.getFFmpegConfig())

	fmt.Println("Done.")
}

func loadConfig() Config {
	configBites, err := ioutil.ReadFile("config.yml")
	if err != nil {

		log.Fatal(err)
	}
	config := Config{}
	err = yaml.Unmarshal(configBites, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return config
}
