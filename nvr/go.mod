module github.com/alexdin/tinygonvr/nvr

go 1.17

replace github.com/alexdin/tinygonvr/ffmpeg => ../ffmpeg

replace github.com/alexdin/tinygonvr/notifyer => ../notifyer

replace github.com/alexdin/tinygonvr/notifyer/telegram => ../notifyer/telegram

replace github.com/alexdin/tinygonvr/alarm => ../alarm

require (
	github.com/alexdin/tinygonvr/alarm v0.0.0-00010101000000-000000000000
	github.com/alexdin/tinygonvr/ffmpeg v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.4.0
)
