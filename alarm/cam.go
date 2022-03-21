package alarm

type Cam struct {
	index int
	ip    string
}

func NewCam(index int, ip string) Cam {
	return Cam{index: index, ip: ip}
}
