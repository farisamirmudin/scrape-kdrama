package model

import (
	"regexp"
	"strings"
)

type Film struct {
	Title string
	Path  string
}

func ExtractTitleAndEpisode(input string) []string {
	matches := regexp.MustCompile(`(.+) Episode (.+)`).FindStringSubmatch(input)
	if len(matches) < 2 {
		return []string{}
	}

	return []string{strings.TrimSpace(matches[1]), strings.TrimSpace(matches[2])}
}
