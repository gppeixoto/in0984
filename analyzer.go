package main

import (
	"fmt"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
)

func analyze(t *TextData, trends []twitter.Trend) (*Response, error) {
	textTokens := getTextTokens(t.Text)
	for _, trend := range trends {
		topicTokens := getTopicTokens(trend.Name)
		if matchedToken, ok := match(textTokens, topicTokens); ok {
			res := &Response{
				Text: fmt.Sprintf(
					"match found %+v matching with \"%v\"",
					trend,
					matchedToken,
				),
			}
			return res, nil
		}
	}
	return nil, fmt.Errorf(fmt.Sprintf("did not find a match for %v", t.Text))
}

func getTopicTokens(str string) (tokens map[string]bool) {
	tokens = make(map[string]bool)
	var ss string
	if str[0] == '#' {
		ss = str[1:len(str)]
	}
	ss = matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	ss = matchAllCap.ReplaceAllString(ss, "${1}_${2}")
	ss = strings.ToLower(ss)

	for _, token := range strings.Split(ss, "_") {
		token = strings.TrimSpace(token)
		tokens[token] = true
	}

	return
}

func getTextTokens(str string) (tokens map[string]bool) {
	tokens = make(map[string]bool)
	for _, token := range strings.Split(str, " ") { // TODO: use proper tokenization
		token = strings.TrimSpace(token)
		if len(token) > 2 {
			tokens[token] = true // instead of splitting on whitespace
		}
	}
	return
}

func match(one, other map[string]bool) (string, bool) {
	for element := range one {
		if other[element] {
			return element, true
		}
	}
	return "", false
}
