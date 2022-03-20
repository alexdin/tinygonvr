module github.com/alexdin/tinygonvr/nvr

go 1.17

replace github.com/alexdin/tinygonvr/ffmpeg => ../ffmpeg

replace github.com/alexdin/tinygonvr/notifyer => ../notifyer

replace github.com/alexdin/tinygonvr/notifyer/telegram => ../notifyer/telegram

require (
	github.com/alexdin/tinygonvr/notifyer v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/alexdin/tinygonvr/notifyer/telegram v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 // indirect
)
