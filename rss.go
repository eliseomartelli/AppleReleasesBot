package main

import "fmt"

type Release struct {
	Title string `xml:"title"`
	Emoji string
	Type  string
}

func (r Release) String() string {
	return fmt.Sprintf("%s %s\n", r.Emoji, r.Title)
}

type RSS struct {
	Releases []Release `xml:"channel>item"`
}
