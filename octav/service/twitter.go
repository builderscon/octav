package service

import (
	"os"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
)

var twitterSvc TwitterSvc
var twitterOnce sync.Once

func Twitter() *TwitterSvc {
	twitterOnce.Do(twitterSvc.Init)
	return &twitterSvc
}

func (v *TwitterSvc) Init() {
	var config oauth2.Config
	var token oauth2.Token

	token.AccessToken = os.Getenv("TWITTER_OAUTH2_ACCESS_TOKEN")
	httpClient := config.Client(oauth2.NoContext, &token)

	twitterSvc.Client = twitter.NewClient(httpClient)
}
