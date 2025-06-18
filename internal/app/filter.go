package app

import (
	"regexp"
	"strings"
)

func CleanMessage(msg string) string {
	re := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	clean := re.ReplaceAllString(strings.ToLower(msg), " ")
	return strings.TrimSpace(clean)
}

func ContainsBadWord(msg string, badWords []string) bool {
	cleanMsg := CleanMessage(msg)
	words := strings.Fields(cleanMsg)

	wordSet := make(map[string]struct{})
	for _, w := range words {
		wordSet[w] = struct{}{}
	}

	for _, bad := range badWords {
		if _, ok := wordSet[bad]; ok {
			return true
		}
	}
	return false
}
