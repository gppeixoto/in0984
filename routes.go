package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	fakeReq       = "{\"text\": \"data\"}"
)

type Server interface {
	Start(port int)
	Close()
}

type server struct {
	server  *echo.Echo
	twitter TwitterTrendsSvc
}

func (s *server) analyzerHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// get data request
		t := new(TextData)
		if err := c.Bind(t); err != nil {
			return c.JSON(http.StatusBadRequest,
				fmt.Sprintf("request should look like %s", fakeReq),
			)
		}
		textTokens := getTextTokens(t.Text)

		// prepare
		trends, err := s.twitter.Trends()
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, "unable to reach twitter")
		}
		for _, trend := range trends {
			topicTokens := getTopicTokens(trend.Name)
			if matchedToken, ok := match(textTokens, topicTokens); ok {
				res := &Response{
					Text: fmt.Sprintf("match found %+v matching with \"%v\"", trend, matchedToken),
				}
				return c.JSON(
					http.StatusOK,
					res,
				)
			}
		}
		return echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("did not find a match for %v", t.Text),
		)
	}
}

func (s *server) registerRoutes() {
	s.server.POST("/text", s.analyzerHandler())
}

func (s *server) Close() {
	s.twitter.Close()
	s.server.Close()
}

func (s *server) Start(port int) {
	s.registerRoutes()
	s.server.Logger.Fatal(s.server.Start(fmt.Sprintf(":%d", port)))
}

func NewServer() Server {
	return &server{
		server:  echo.New(),
		twitter: NewTwitterTrendsSvc(23424768), // Brazil WOEID
	}
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
		tokens[token] = true
	}

	return
}

func getTextTokens(str string) (tokens map[string]bool) {
	tokens = make(map[string]bool)
	for _, token := range strings.Split(str, " ") { // TODO: use proper tokenization
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
