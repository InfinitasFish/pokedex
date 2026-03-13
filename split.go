package main

import "strings"


func cleanInput(text string) []string {
	tokens := strings.Fields(text)
	lowerTokens := make([]string, 0, len(tokens))
	for _, token := range tokens {
		lowerTokens = append(lowerTokens, strings.ToLower(token))
	}
	return lowerTokens
}

