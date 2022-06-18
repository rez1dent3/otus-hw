package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var onlyWord = regexp.MustCompile("[A-Za-zА-ЯЁа-яё][A-Za-zА-ЯЁа-яё-]*")

func Top10(input string) []string {
	text := strings.ToLower(input)
	wordCount := make(map[string]int)
	for _, word := range onlyWord.FindAllString(text, -1) {
		wordCount[word]++
	}

	if len(wordCount) == 0 {
		return nil
	}

	i := 0
	words := make([]string, len(wordCount))
	for key := range wordCount {
		words[i] = key
		i++
	}

	sort.Slice(words, func(i, j int) bool {
		if wordCount[words[i]] == wordCount[words[j]] {
			return words[i] < words[j]
		}

		return wordCount[words[i]] > wordCount[words[j]]
	})

	length := len(words)
	if length > 10 {
		length = 10
	}

	return words[:length]
}
