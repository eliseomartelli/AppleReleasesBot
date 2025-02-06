package main

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

func main() {
	config := LoadConfig()
	constants := SetupConstants()

	interval := flag.Duration("interval", 0, "Polling interval for continuous checking (e.g., 5m, 1h)")
	flag.Parse()

	if *interval <= 0 {
		fmt.Println("Please specify a valid polling interval.")
		return
	}
	ticker := time.NewTicker(*interval)

	fmt.Printf("Starting task every %s...\n", *interval)

	for range ticker.C {
		checkAndNotify(config, constants)
	}
}

func checkAndNotify(config *config, constants appConstants) {
	if releases, err := fetchReleases(config.RSSFeedURL); err != nil {
		log.Fatalf("Error fetching releases: %v", err)
	} else {
		processed := processReleases(releases, constants.NameToEmoji)
		releaseText := generateReleaseText(processed)
		fmt.Print(releaseText)
		if shouldPostUpdate(config.CacheFilePath, releaseText) {
			if err := sendTelegramNotification(config, releaseText); err != nil {
				log.Fatalf("Error sending notification: %v", err)
			}
			updateCache(config.CacheFilePath, releaseText)
		}
	}
}

func fetchReleases(url string) ([]Release, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", SetupConstants().UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	var rss RSS
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return nil, fmt.Errorf("failed to decode RSS feed: %w", err)
	}

	return rss.Releases, nil
}

func processReleases(releases []Release, emojiMap map[string]string) []Release {
	var filtered []Release
	for _, r := range releases {
		if osType := SetupConstants().OSPattern.FindString(r.Title); osType != "" {
			r.Type = osType
			r.Emoji = emojiMap[osType]
			filtered = append(filtered, r)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Title < filtered[j].Title
	})

	return filtered
}

func generateReleaseText(releases []Release) string {
	var buf bytes.Buffer
	for _, r := range releases {
		buf.WriteString(r.String())
	}
	return buf.String()
}

func shouldPostUpdate(cachePath, newContent string) bool {
	oldContent, _ := os.ReadFile(cachePath)
	return md5.Sum([]byte(newContent)) != md5.Sum(oldContent)
}

func updateCache(cachePath, content string) {
	if err := os.WriteFile(cachePath, []byte(content), 0o600); err != nil {
		log.Printf("Warning: Failed to update cache file: %v", err)
	}
}
