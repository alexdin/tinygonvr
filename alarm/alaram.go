package alarm

type Alarm struct {
	cams    []Cam
	channel chan int
}

var alarm Alarm

func Boot(config []Cam) {
	alarm = Alarm{
		cams:    config,
		channel: make(chan int, len(config)),
	}
}

func (a *Alarm) GetReadChannel() <-chan int {
	return a.channel
}

func (a *Alarm) sendAlarmToChannel(index int) {
	a.channel <- index
}

func GetAlarm() Alarm {
	return alarm
}
