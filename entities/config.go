package entities

type Config struct {
	Token          string
	ChatId         int64  `split_words:"true"`
	WelcomeSticker string `split_words:"true"`
	Debug          bool
}
