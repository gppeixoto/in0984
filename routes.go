package main

import (
	"fmt"
	"math/rand"
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

// Server is the main interface to implement a server
type Server interface {
	Start(port int)
	Close()
}

type server struct {
	server    *echo.Echo
	twitter   TwitterTrendsSvc
	sentiment SentSvc
}

func (s *server) analyzerHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// try twitter connection
		trends, err := s.twitter.Trends()
		if err != nil {
			echo.NewHTTPError(
				http.StatusInternalServerError,
				"unable to reach twitter")
		}
		// get data request
		t := new(TextData)
		if err := c.Bind(t); err != nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				fmt.Sprintf("request should look like %s", fakeReq))
		}
		match, err := matchAndGetQuery(t, trends)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		tweets, err := s.twitter.TweetsFor(match.query)
		if err != nil {
			//TODO: check error msg
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		// pass tweets to sentiment analyzer
		score, magnitude := s.sentiment.Score(tweets)
		r := &Response{
			Name:         match.name,
			Score:        score,
			TweetVolume:  match.volume,
			Magnitude:    magnitude,
			SampleTweets: sampleTweets(tweets, t.TweetsSampleSize)}
		return c.JSON(http.StatusOK, r)
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

// NewServer creates a new app server
func NewServer() Server {
	return &server{
		server:    echo.New(),
		twitter:   NewTwitterTrendsSvc(23424768), // Brazil WOEID
		sentiment: NewSentSvc(),
	}
}

func sampleTweets(tweets []string, n int) []string {
	n = min(n, 10)
	tmp := make([]string, len(tweets))
	copy(tmp, tweets)
	for i := range tmp {
		j := rand.Intn(i + 1)
		tmp[i], tmp[j] = tmp[j], tmp[i]
	}
	var resp []string
	for i := 0; i < len(tmp) && len(resp) < n; i++ {
		if !strings.HasPrefix(tmp[i], "RT") {
			resp = append(resp, tmp[i])
		}
	}
	return resp
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
