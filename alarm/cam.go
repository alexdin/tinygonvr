package alarm

type Cam struct {
	index   int
	ip      string
	channel chan int
}

func NewCam(index int, ip string) Cam {
	return Cam{index: index, ip: ip, channel: make(chan int)}
}
