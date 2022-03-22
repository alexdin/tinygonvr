package main

import (
	"github.com/alexdin/tinygonvr/alarm"
	"github.com/alexdin/tinygonvr/ffmpeg"
	"net"
	"net/url"
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

func (c *Config) getAlertConfig() []alarm.Cam {
	var slice []alarm.Cam
	for index, val := range c.Cams {
		slice = append(slice, alarm.NewCam(index, getIpFromUrl(val.Url)))
	}
	return slice
}

func (c *Config) getFFmpegConfig() []ffmpeg.Stream {
	var slice []ffmpeg.Stream
	for _, val := range c.Cams {
		slice = append(slice, ffmpeg.Stream{Url: val.Url, CamName: val.Name})
	}
	return slice
}

func getIpFromUrl(parse string) string {
	u, err := url.Parse(parse)
	if err != nil {
		panic(err)
	}
	host, _, _ := net.SplitHostPort(u.Host)
	return host
}
