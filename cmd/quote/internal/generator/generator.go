package generator

import (
	"math/rand"
	"time"
)

var quotes = []string{
	"don't put your hand in boiling water",
	"do not breathe under the water",
	"breathing will help you live",
	"don't look up and spit",
	"in case of fire, exit building before tweeting about it",
	"don't take advice from posters",
	"don't eat yellow snow",
	"don't swim in waters inhabited by large alligators",
	"a day without sunshine is like, you know, night",
	"there isn't really a Nigerian Prince who wants to transfer money to you",
	"never buy a car you can’t push",
	"only ninja can sneak upon other ninja",
}
var randomSource *rand.Rand

func init() {
	randomSource = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
}

func Get() string {
	return quotes[randomSource.Intn(len(quotes))]
}
