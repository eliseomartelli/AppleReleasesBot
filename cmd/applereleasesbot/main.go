package main

import (
	"applereleases/internal/config"
	"applereleases/internal/model"
	"applereleases/internal/service/apple"
	"applereleases/internal/service/telegram"
	"applereleases/internal/state"
	"bytes"
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	appleService := apple.NewService()
	telegramService := telegram.NewService(cfg.TelegramBotToken, cfg.TelegramChatID)
	stateManager := state.NewManager(cfg.CacheFilePath)

	interval := flag.Duration("interval", 0, "Polling interval for continuous checking (e.g., 5m, 1h). Set to 0 for one-shot run.")
	flag.Parse()

	if *interval > 0 {
		log.Printf("Starting task every %s...\n", *interval)
		ticker := time.NewTicker(*interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := run(cfg, appleService, telegramService, stateManager); err != nil {
				log.Printf("Error running task: %v", err)
			}
		}
	} else {
		log.Println("Running one-shot task...")
		if err := run(cfg, appleService, telegramService, stateManager); err != nil {
			log.Fatalf("Error running task: %v", err)
		}
	}
}

func run(cfg *config.Config, appleService *apple.Service, telegramService *telegram.Service, stateManager *state.Manager) error {
	releases, err := appleService.FetchReleases(cfg.RSSFeedURL)
	if err != nil {
		return fmt.Errorf("error fetching releases: %w", err)
	}

	processed := appleService.ProcessReleases(releases)
	releaseText := generateReleaseText(processed)

	// Original logic printed the text to stdout
	fmt.Print(releaseText)

	if stateManager.ShouldPostUpdate(releaseText) {
		log.Println("New releases found, sending notification...")
		if err := telegramService.SendNotification(releaseText); err != nil {
			return fmt.Errorf("error sending notification: %w", err)
		}
		if err := stateManager.UpdateCache(releaseText); err != nil {
			log.Printf("Warning: failed to update cache: %v", err)
		}
	}

	return nil
}

func generateReleaseText(releases []model.Release) string {
	var buf bytes.Buffer
	for _, r := range releases {
		buf.WriteString(r.String())
	}
	buf.WriteString("\nhttps://support.apple.com/en-us/100100")
	return buf.String()
}
