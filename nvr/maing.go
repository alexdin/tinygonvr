package main

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
// #cgo CFLAGS: -std=c11 -g
// #include "ffmpeg.h"
import "C"
import (
	"fmt"
	notifyer "github.com/alexdin/tinygonvr/notifyer"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Camera struct {
	Url  string `yaml:"url"`
	Name string `yaml:"name"`
}

type Config struct {
	Debug  bool
	Cams   []Camera `yaml:"cams"`
	Notify Notify   `yaml:"notify"`
}

type Notify struct {
	BotToken  string `yaml:"botToken"`
	ChannelId int    `yaml:"channelId"`
}

func main() {

	config := loadConfig()

	for _, cam := range config.Cams {
		fmt.Println(cam)
		/*	stream := ffmpeg.Stream{Url: cam.Url, CamName: cam.Name}
			stream.Open()
			stream.Screen()
			stream.Close()*/

	}

	fmt.Println("Done")
	notifyer.Boot(
		notifyer.Notify{
			BotType:   notifyer.ChannelTelegram,
			BotToken:  config.Notify.BotToken,
			ChannelId: config.Notify.ChannelId,
		},
	)
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
