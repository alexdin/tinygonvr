package main

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
// #cgo CFLAGS: -std=c11 -g
// #include "ffmpeg.h"
import "C"
import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/alexdin/tinygonvr/ffmpeg"

	"gopkg.in/yaml.v2"
)

type Camera struct {
	Url  string `yaml:"url"`
	Name string `yaml:"name"`
}

type Config struct {
	Debug bool
	Cams  []Camera `yaml:"cams"`
}

func main() {

	config := loadConfig()

	for _, cam := range config.Cams {
		stream := ffmpeg.Stream{Url: cam.Url, CamName: cam.Name}
		stream.Open()
		stream.Screen()
		stream.Close()

	}

	fmt.Println("Done")
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
