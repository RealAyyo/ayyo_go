package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reg = regexp.MustCompile(`[\p{L}0-9-]+(?:['!.,]?[\p{L}0-9-]+)*`)

func Top10(str string) []string {
	splitString := reg.FindAllString(str, -1)
	wordFrequency := make(map[string]int)

	for _, word := range splitString {
		wordLower := strings.ToLower(word)
		if word == "" || word == "-" {
			continue
		}
		wordFrequency[wordLower]++
	}

	wordResult := make([]string, 0, len(wordFrequency))

	for k := range wordFrequency {
		wordResult = append(wordResult, k)
	}

	sort.Slice(wordResult, func(i, j int) bool {
		if wordFrequency[wordResult[i]] != wordFrequency[wordResult[j]] {
			return wordFrequency[wordResult[i]] > wordFrequency[wordResult[j]]
		}
		return wordResult[i] < wordResult[j]
	})

	if len(wordResult) > 10 {
		return wordResult[:10]
	}
	return wordResult
}
