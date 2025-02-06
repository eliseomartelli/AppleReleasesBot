package main

import (
	"log"
	"os"
)

type config struct {
	TelegramBotToken string
	TelegramChatID   string
	RSSFeedURL       string
	CacheFilePath    string
}

func LoadConfig() *config {
	return &config{
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:   getEnv("TELEGRAM_CHAT_ID", ""),
		RSSFeedURL:       getEnv("RSS_FEED_URL", "https://developer.apple.com/news/releases/rss/releases.rss"),
		CacheFilePath:    getEnv("CACHE_FILE_PATH", "./newscache"),
	}
}

func getEnv(key, defaulValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if defaulValue == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return defaulValue
}
