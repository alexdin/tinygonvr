package alarm

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Alarm struct {
	cams []Cam
}

type Message struct {
	Address   string `json:"Address"`
	Channel   int    `json:"Channel"`
	Descrip   string `json:"Descrip"`
	Event     string `json:"Event"`
	Serialid  string `json:"SerialID"`
	Starttime string `json:"StartTime"`
	Status    string `json:"Status"`
	Type      string `json:"Type"`
}

var alarm Alarm

func Boot(config []Cam) {
	alarm = Alarm{
		cams: config,
	}
	// run alarm server
	//serverBoot()

	// send bytes
	s := "{\"Address\":\"0x671FA8C0\",\"Channel\":0,\"Descrip\":\"\",\"Event\":\"HumanDetect\",\"SerialID\":\"783ae37bbef6379d\",\"StartTime\":\"2022-03-21 21:19:34\",\"Status\":\"Start\",\"Type\":\"Alarm\",\"ipAddr\":\"192.168.31.103\"}\n"
	handleAlert([]byte(s))
}

func (a *Alarm) GetCamChannelByIndex(index int) <-chan int {
	for _, v := range a.cams {
		if v.index == index {
			return v.channel
		}
	}
	return nil
}

func GetAlarm() Alarm {
	return alarm
}

func serverBoot() {
	processesWaitGroup := sync.WaitGroup{}

	server := Server{
		Port:           "15002",
		WaitGroup:      &processesWaitGroup,
		MessageHandler: handleAlert,
	}
	server.Start()

	processesWaitGroup.Wait()
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

func handleAlert(data []byte) {
	message := Message{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		log.Fatal(err)
		return
	}
	if message.Status == "Start" && message.Event == "HumanDetect" && message.Type == "Alarm" {
		for _, val := range alarm.cams {
			if val.ip == hexIpToCIDR(message.Address) {
				val.channel <- val.index
				return
			}
		}
	}
}
