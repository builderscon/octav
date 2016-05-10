package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"syscall"

	"github.com/builderscon/octav/worker/deploybot/slacksub"
	"github.com/lestrrat/go-pdebug"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
	"google.golang.org/cloud/pubsub"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	var authtokenf string
	var email string
	var fifopath string
	var projectID string
	var slackgw string
	var topic string
	var zone string
	flag.StringVar(&authtokenf, "authtokenfile", "", "File containing token used to authentication when posting")
	flag.StringVar(&email, "email", "", "email ID to use for acme protocol")
	flag.StringVar(&fifopath, "fifopath", "", "path to where tls requests willbe pushed to")
	flag.StringVar(&projectID, "project_id", "", "project ID to use")
	flag.StringVar(&slackgw, "slackgw", "http://slackgw:4979", "slack gateway url")
	flag.StringVar(&topic, "topic", "slackgw-url", "topic name to subscribe to")
	flag.StringVar(&zone, "zone", "", "DNS zone to update")
	flag.Parse()

	var authtoken string
	if authtokenf != "" {
		buf, err := ioutil.ReadFile(authtokenf)
		if err != nil {
			fmt.Printf("Failed to open file '%s': %s", authtokenf, err)
			return 1
		}
		authtoken = string(buf)
	}

	if fifopath == "" {
		fmt.Printf("fifopath is required")
		return 1
	}

	if projectID == "" {
		fmt.Printf("projectID is required")
		return 1
	}

	pcl, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		fmt.Printf("failed to create pubsub client: %s", err)
		return 1
	}

	bot := New(pcl, topic, slackgw, authtoken, fifopath)
	bot.Run()

	return 0
}

type Bot struct {
	*slacksub.Subscriber
	fifopath string
}

type deployctx struct {
	bot *Bot
	msg *slack.MessageEvent
}

type ingressctx struct {
	bot *Bot
	msg *slack.MessageEvent
}

func New(cl *pubsub.Client, topic, slackgwURL, token, fifopath string) *Bot {
	bot := &Bot{
		Subscriber: slacksub.New(cl, topic, slackgwURL, token),
		fifopath:   fifopath,
	}
	bot.Subscriber.MessageCallback = bot.processMessageEvent
	return bot
}

func (b *Bot) Fifo() (f *os.File, err error) {
	path := b.fifopath
	if pdebug.Enabled {
		g := pdebug.Marker("b.Fifo %s", path).BindError(&err)
		defer g.End()
	}

	if _, err := os.Stat(path); err != nil { // doesn't exist
		if err := syscall.Mknod(path, syscall.S_IFIFO|0666, 0); err != nil {
			// Failed to create... timing problem?
			if _, err := os.Stat(path); err != nil {
				// Hmm, weird. bail
				return nil, err
			}
		}
	}

	fh, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return fh, nil
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

func (ctx *ingressctx) Reply(s string) error {
	return ctx.bot.Subscriber.Reply(ctx.msg.Channel, s)
}

type replier interface {
	Reply(string) error
}

func (b *Bot) handleHelpCmd(ctx replier) error {
	return ctx.Reply(`botname ingress create <domain> [key=value ...]
botname ingress delete <domain>
botname ingress get <domain>
botname ingress list`)
}

var spacesRx = regexp.MustCompile(`\s+`)

func (b *Bot) processMessageEvent(ev *slack.MessageEvent) error {
	if pdebug.Enabled {
		g := pdebug.Marker("b.processMessageEvent")
		defer g.End()
	}

	cmd := spacesRx.Split(strings.TrimSpace(ev.Text), -1)
	if len(cmd) < 3 {
		return nil
	}

	switch cmd[1] {
	case "ingress":
		ctx := ingressctx{bot: b, msg: ev}
		return b.handleIngressCmd(&ctx, cmd[2:])
	case "deploy":
		ctx := deployctx{bot: b, msg: ev}
		return b.handleDeployCmd(&ctx, cmd[2:])
	}
	return nil
}

type deployargs struct {
	Target  string            `json:"target"`
	Channel string            `json:"channel"`
	Name    string            `json:"name"`
	Mode    string            `json:"mode"`
	Args    map[string]string `json:"args"`
}

// @gkebot ingress create <fqdn> [key=value ...]
// @gkebot ingress delete <fqdn>
// @gkebot ingress get <fqdn>
func (b *Bot) handleIngressCmd(ctx *ingressctx, cmd []string) error {
	if len(cmd) < 1 {
		return b.handleHelpCmd(ctx)
	}

	args := deployargs{
		Target:  "ingress",
		Channel: ctx.msg.Channel,
	}
	if cmd[0] != "list" {
		if len(cmd) < 2 {
			return b.handleHelpCmd(ctx)
		}

		var hostname string
		sl, err := parseSlackLink(cmd[1])
		if err == nil {
			hostname = sl.Text
		} else {
			hostname = cmd[1]
		}

		args.Name = hostname
	}

	switch cmd[0] {
	case "list":
		args.Mode = "list"
	case "get":
		args.Mode = "get"
	case "delete":
		args.Mode = "delete"
	case "create":
		args.Mode = "create"

		if len(cmd) > 3 {
			args.Args = map[string]string{}
			for _, carg := range cmd[2:] {
				pair := strings.SplitN(carg, "=", 2)
				if len(pair) != 2 {
					return ctx.Reply("expected key=value pairs but got '" + carg + "'")
				}
				args.Args[pair[0]] = pair[1]
			}
		}
	default:
		return b.handleHelpCmd(ctx)
	}

	// There's nothing to be done using GCP stuff here, so pass it
	// to the worker that's linked with k8s libraries. Library
	// version incompatibilities suck.
	fifo, err := b.Fifo()
	if err != nil {
		return err
	}
	defer fifo.Close()

	// Write to local fifo
	if err := json.NewEncoder(fifo).Encode(args); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleDeployCmd(ctx *deployctx, cmd []string) error {
	return nil
}
