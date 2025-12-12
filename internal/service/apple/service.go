package apple

import (
	"applereleases/internal/model"
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"sort"
)

type Service struct {
	client    *http.Client
	osPattern *regexp.Regexp
	emojiMap  map[string]string
	userAgent string
}

func NewService() *Service {
	return &Service{
		client:    http.DefaultClient,
		osPattern: regexp.MustCompile(`(macOS|iOS|iPadOS|watchOS|visionOS|tvOS)`),
		emojiMap: map[string]string{
			"macOS":    "💻",
			"iOS":      "📱",
			"iPadOS":   "📱",
			"watchOS":  "⌚️",
			"visionOS": "🥽",
			"tvOS":     "📺",
		},
		userAgent: "AppleReleasesBot/1.0 (+https://github.com/eliseomartelli/AppleReleasesBot)",
	}
}

func (s *Service) FetchReleases(url string) ([]model.Release, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", s.userAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	var rss model.RSS
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return nil, fmt.Errorf("failed to decode RSS feed: %w", err)
	}

	return rss.Releases, nil
}

func (s *Service) ProcessReleases(releases []model.Release) []model.Release {
	var filtered []model.Release
	for _, r := range releases {
		if osType := s.osPattern.FindString(r.Title); osType != "" {
			r.Type = osType
			r.Emoji = s.emojiMap[osType]
			filtered = append(filtered, r)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Title < filtered[j].Title
	})

	return filtered
}
