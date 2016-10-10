package service

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"sync"
	"unicode/utf8"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"
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

func (v *TwitterSvc) TweetAsConference(confID, tweet string) error {
	// Post to twitter, but we can only do so if we have a valid
	// credential information. This is stored in Google storage
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// twitter credentials
	credentialsKey := "conferences/" + confID + "/credentials/twitter"
	var credentialsBuf bytes.Buffer
	if err := CredentialStorage.Download(ctx, credentialsKey, &credentialsBuf); err != nil {
		return errors.Wrap(err, "failed to download twitter credentials")
	}

	// ...and they are in JSON
	var creds struct {
		AccessToken  string `json:"access_token"`
		AccessSecret string `json:"access_scret"`
	}

	if err := json.Unmarshal(credentialsBuf.Bytes(), &creds); err != nil {
		return errors.Wrap(err, "failed to unmarshal twitter credentials")
	}

	// Consumer key and secret are from env vars
	consumerKey := os.Getenv("TWITTER_OAUTH1_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_OAUTH1_CONSUMER_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	if utf8.RuneCountInString(tweet) > 140 {
		var truncated bytes.Buffer
		for i := 1; i < 140; i++ { // 139 runes
			r, n := utf8.DecodeRuneInString(tweet)
			if r == utf8.RuneError {
				break
			}
			tweet = tweet[n:]
			truncated.WriteRune(r)
		}
		truncated.WriteRune('…')
		tweet = truncated.String()
	}

	client.Statuses.Update(tweet, nil)
	return nil
}
