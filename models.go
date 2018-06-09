package main

type TextData struct {
	Text string `json:"text"`
}

type Response struct {
	Name        string  `json:"name"`
	Score       float64 `json:"sentiment_score"`
	Magnitude   float64 `json:"sentiment_strength"`
	TweetVolume int64   `json:"tweet_volume"`
}

type Match struct {
	query  string
	name   string
	volume int64
}
