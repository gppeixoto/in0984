package main

import (
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// TwitterTrendsSvc accesses twitter
type TwitterTrendsSvc interface {
	Trends() ([]twitter.Trend, error)
	Close()
}

type twitterTrends struct {
	client *twitter.Client
	woeid  int64
}

func (tt *twitterTrends) Close() {}

func (tt *twitterTrends) Trends() ([]twitter.Trend, error) {
	ts, _, err := tt.client.Trends.Place(tt.woeid, nil)
	if err != nil {
		return nil, err
	}
	var trends []twitter.Trend
	for _, xyz := range ts {
		for _, trend := range xyz.Trends {
			trends = append(trends, trend)
		}
	}
	return trends, nil
}

// NewTwitterTrendsSvc creates a new TwitterTrendsSvc
func NewTwitterTrendsSvc(woeid int64) TwitterTrendsSvc {
	return &twitterTrends{
		client: newClient(),
		woeid:  woeid,
	}
}

func newClient() *twitter.Client {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	client := twitter.NewClient(config.Client(oauth1.NoContext, token))
	return client
}

type fakeTwitter struct{}

func (ft *fakeTwitter) Trends() ([]twitter.Trend, error) {
	return []twitter.Trend{
		twitter.Trend{
			Name:        "chewbacca",
			URL:         "chewbac.ca",
			TweetVolume: 42,
		},
	}, nil
}

func (ft *fakeTwitter) Close() {}

// FakeTwitterTrendsSvc does the boomshakalaka
func FakeTwitterTrendsSvc(woeid int64) TwitterTrendsSvc {
	return &fakeTwitter{}
}
