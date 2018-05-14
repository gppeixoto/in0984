package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}

	twitter := NewTwitterTrendsSvc(23424768) // Brazil WOEID
	trends, err := twitter.Trends()
	if err != nil {
		panic(err)
	}

	for _, trend := range trends {
		fmt.Printf("%+v\n", trend)
	}
}
