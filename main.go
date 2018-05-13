package main

import "fmt"

func main() {
	var twitter TwitterTrendsSvc
	twitter = FakeTwitterTrendsSvc(0)
	trends, err := twitter.Trends()
	if err != nil {
		panic(err)
	}
	for _, trend := range trends {
		fmt.Printf("%+v", trend)
	}
}
