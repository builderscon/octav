package slackbot

import (
	"errors"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/lestrrat/go-cloud-acmeagent"
	"github.com/lestrrat/go-pdebug"
	"github.com/nlopes/slack"
)

type rtmctx struct {
	RTM     *slack.RTM
	Message *slack.MessageEvent
}

func (ctx *rtmctx) Reply(txt string) {
	ctx.RTM.SendMessage(ctx.RTM.NewOutgoingMessage(txt, ctx.Message.Channel))
}

func StartRTM(done chan struct{}) {
	defer close(done)
	if pdebug.Enabled {
		g := pdebug.Marker("StartRTM")
		defer g.End()
	}

	rtm := slackClient.NewRTM()
	go rtm.ManageConnection()

	sigCh := make(chan os.Signal, 265)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for loop := true; loop; {
		select {
		case msg := <-rtm.IncomingEvents:
			if err := handleMessage(rtm, msg); err != nil {
				if pdebug.Enabled {
					pdebug.Printf("handleMessage: %s", err)
				}
				loop = false
			}
		case <-sigCh:
			loop = false
		case <-done:
			loop = false
		}
	}
}

type SlackLink struct {
	Text string
	URL  string
}

func parseSlackLink(s string) (*SlackLink, error) {
	if len(s) == 0 || s[0] != '<' {
		return nil, errors.New("not a link")
	}
	sl := &SlackLink{}
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '|':
			sl.Text = s[1:i]
		case '>':
			if l := len(sl.Text); l > 0 {
				sl.URL = sl.Text
				sl.Text = s[len(sl.Text)+2 : i]
			} else {
				sl.Text = s[1:i]
			}
			return sl, nil
		}
	}

	return nil, errors.New("not a link")
}

var spacesRx = regexp.MustCompile(`\s+`)

func handleMessage(rtm *slack.RTM, msg slack.RTMEvent) (err error) {
	switch msg.Data.(type) {
	case *slack.RTMError:
		return msg.Data.(*slack.RTMError)
	case *slack.InvalidAuthEvent:
		return errors.New("invalid auth")
	case *slack.MessageEvent:
		sm := msg.Data.(*slack.MessageEvent)

		cmd := spacesRx.Split(strings.TrimSpace(sm.Text), -1)
		if len(cmd) < 3 {
			return nil
		}

		sl, err := parseSlackLink(cmd[0])
		if err != nil || sl.Text != "@"+slackUser {
			return nil
		}

		if cmd[1] != "acme" {
			return nil
		}
		ctx := rtmctx{RTM: rtm, Message: sm}
		handleLetsEncryptCmd(&ctx, cmd[2:])
	}

	return nil
}

func handleLetsEncryptCmd(ctx *rtmctx, cmd []string) {
	switch cmd[0] {
	case "authz":
		if len(cmd) < 2 {
			return
		}
		handleAuthzCmd(ctx, cmd[1:])
	case "cert":
		if len(cmd) < 2 {
			return
		}
		handleCertCmd(ctx, cmd[1:])
	case "upload":
		if len(cmd) < 2 {
			return
		}
		handleUploadCmd(ctx, cmd[1:])
	}
}

func handleAuthzCmd(ctx *rtmctx, cmd []string) {
	if len(cmd) < 1 {
		return
	}

	sl, err := parseSlackLink(cmd[0])
	if err != nil {
		return
	}

	domain := sl.Text

	ctx.Reply(":white_check_mark: Authorizing *" + domain + "*")

	var authz acmeagent.Authorization
	if err := acmeStateStore.LoadAuthorization(domain, &authz); err != nil {
		ctx.Reply(":white_check_mark: Authorization for domain not found in storage.")
	} else {
		if authz.IsExpired() {
			ctx.Reply(":exclamation: Authorization expired, going to run authorization again")
		} else {
			ctx.Reply(":exclamation: Authorization already exists. Run `acme cert` to issue certificates for this domain")
			return
		}
	}

	ctx.Reply(":white_check_mark: Running authorization (this may take a few minutes)")
	// Do this in a goroutine so we don't block from doing other things
	go func() {
		if err := acmeAgent.AuthorizeForDomain(domain); err != nil {
			ctx.Reply(":exclamation: Authorization failed: " + err.Error())
			return
		}
		ctx.Reply(":tada: Authorization for domain *" + domain + "* complete")
	}()
}

func handleCertCmd(ctx *rtmctx, cmd []string) {
	if len(cmd) < 1 {
		return
	}

	sl, err := parseSlackLink(cmd[0])
	if err != nil {
		return
	}

	domain := sl.Text
	ctx.Reply(":white_check_mark: Issueing certificates for *" + domain + "*")

	cert, err := acmeStateStore.LoadCert(domain)
	if err != nil {
		ctx.Reply(":white_check_mark: Certificates for domain not found in storage.")
	} else {
		if time.Now().After(cert.NotAfter) {
			ctx.Reply(":exclamation: Certificate expired, going to issue it again")
		} else {
			ctx.Reply(":exclamation: Certificate already exists. Run `acme upload` to upload the certificate")
			return
		}
	}

	// run handleAuthzCmd to make sure that the authorization is there
	handleAuthzCmd(ctx, cmd)

	ctx.Reply(":white_check_mark: Fetching certificates")
	// Do this in a goroutine so we don't block from doing other things
	go func() {
		if err := acmeAgent.IssueCertificate(domain, nil, false); err != nil {
			ctx.Reply(":exclamation: Failed to fetch certificates: " + err.Error())
			return
		}
		ctx.Reply(":tada: Authorization for domain *" + domain + "* complete")
	}()
}

func handleUploadCmd(ctx *rtmctx, cmd []string) {
	if len(cmd) < 1 {
		return
	}

	sl, err := parseSlackLink(cmd[0])
	if err != nil {
		return
	}

	domain := sl.Text
	ctx.Reply(":white_check_mark: Uploading certificates for *" + domain + "*")

	name, err := acmeAgent.UploadCertificate(domain)
	if err != nil {
		ctx.Reply(":exclamation: Failed to upload certificates: " + err.Error())
		return
	}

	ctx.Reply(":tada: Certificates uploaded as *" +name+"*")
}
