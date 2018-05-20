package main

type TextData struct {
	Text string `json:"text"`
}

type Response struct {
	Name        string  `json:"name"`
	Score       float32 `json:"score"`
	TweetVolume int64   `json:"tweet_volume"`
}

type Match struct {
	query  string
	name   string
	volume int64
}
