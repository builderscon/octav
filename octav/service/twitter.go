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
	twitterSvc.Client = NewTwitterClientFromToken(os.Getenv("TWITTER_OAUTH2_ACCESS_TOKEN"))
}

func NewTwitterClientFromToken(s string) *twitter.Client {
	var config oauth2.Config
	var token oauth2.Token

	token.AccessToken = s
	httpClient := config.Client(oauth2.NoContext, &token)
	return twitter.NewClient(httpClient)
}

