package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type Server interface {
	Start(port int)
	Close()
}

type server struct {
	server *echo.Echo
}

type TextData struct {
	Text string `json:"text"`
}

func (s *server) analyzerHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		t := new(TextData)
		if err := c.Bind(t); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, t)
	}
}

func (s *server) registerRoutes() {
	s.server.POST("/text", s.analyzerHandler())
}

func (s *server) Close() {
	s.server.Close()
}

func (s *server) Start(port int) {
	s.registerRoutes()
	s.server.Logger.Fatal(s.server.Start(fmt.Sprintf(":%d", port)))
}

func NewServer() Server {
	return &server{
		server: echo.New(),
	}
}
