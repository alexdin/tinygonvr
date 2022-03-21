package main

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
// #cgo CFLAGS: -std=c11 -g
// #include "../ffmpeg/ffmpeg.h"
import "C"
import (
	"fmt"
	"github.com/alexdin/tinygonvr/alarm"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func main() {

	// load config data
	config := loadConfig()
	// boot config alert monitoring
	alarm.Boot(config.getAlertConfig())
	for _, cam := range config.Cams {
		fmt.Println(cam)
		/*	stream := ffmpeg.Stream{Url: cam.Url, CamName: cam.Name}
			stream.Open()
			stream.Screen()
			stream.Close()*/

	}

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
