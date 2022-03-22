package main

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
// #cgo CFLAGS: -std=c11 -g
// #include "../ffmpeg/ffmpeg.h"
import "C"
import (
	"fmt"
	"github.com/alexdin/tinygonvr/alarm"
	"github.com/alexdin/tinygonvr/ffmpeg"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func main() {

	// load config data
	config := loadConfig()
	// boot config alert monitoring
	go alarm.Boot(config.getAlertConfig())
	ffmpeg.Boot(config.getFFmpegConfig())
	fmt.Println("Done")
	//notifyer.Boot(
	//	notifyer.Notify{
	//		BotType:   notifyer.ChannelTelegram,
	//		BotToken:  config.Notify.BotToken,
	//		ChannelId: config.Notify.ChannelId,
	//	},
	//)
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
