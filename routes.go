package main

import (
	"fmt"
	"net/http"
	"regexp"

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
		r, err := analyze(t, trends)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
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

func NewServer() Server {
	return &server{
		server:  echo.New(),
		twitter: NewTwitterTrendsSvc(23424768), // Brazil WOEID
	}
}
