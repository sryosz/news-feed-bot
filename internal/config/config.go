package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env                  string        `yaml:"env" env-default:"local"`
	TelegramBotToken     string        `yaml:"telegram_bot_token" env-required:"true"`
	TelegramChannelID    int64         `yaml:"telegram_channel_id"  env-required:"true"`
	DatabaseDSN          string        `yaml:"database_dsn" env-default:"postgres://postgres:postgres@localhost:5432/news_feed_bot?sslmode=disable"`
	FetchInterval        time.Duration `yaml:"fetch_interval" env-default:"10m"`
	NotificationInterval time.Duration `yaml:"notification_interval" env-default:"10m"`
	FilterKeywords       []string      `yaml:"filter_keywords" `
	OpenAIKey            string        `yaml:"openai_key"`
	OpenAIPrompt         string        `yaml:"openai_prompt"`
	OpenAIModel          string        `yaml:"openai_model" env-default:"gpt-3.5-turbo"`
}

func MustLoad() *Config {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
