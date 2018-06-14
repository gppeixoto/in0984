package main

// TextData defines model for api request
type TextData struct {
	Text             string `json:"text"`
	TweetsSampleSize int    `json:"tweets_sample_size"`
}

// Response defines API's response
type Response struct {
	Name         string   `json:"name"`
	Score        float64  `json:"sentiment_score"`
	Magnitude    float64  `json:"sentiment_strength"`
	TweetVolume  int64    `json:"tweet_volume"`
	SampleTweets []string `json:"sample_tweets"`
}

type trendingTopicMatch struct {
	query  string
	name   string
	volume int64
}
