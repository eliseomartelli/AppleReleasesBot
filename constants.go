package main

import "regexp"

type appConstants struct {
	OSPattern   *regexp.Regexp
	NameToEmoji map[string]string
	UserAgent   string
}

func SetupConstants() appConstants {
	osPattern := regexp.MustCompile(`(macOS|iOS|iPadOS|watchOS|visionOS|tvOS)`)

	return appConstants{
		OSPattern: osPattern,
		NameToEmoji: map[string]string{
			"macOS":    "💻",
			"iOS":      "📱",
			"iPadOS":   "📱",
			"watchOS":  "⌚️",
			"visionOS": "🥽",
			"tvOS":     "📺",
		},
		UserAgent: "AppleReleasesBot/1.0 (+https://github.com/eliseomartelli/AppleReleasesBot)",
	}
}
