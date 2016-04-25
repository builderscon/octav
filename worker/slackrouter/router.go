package main

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/builderscon/octav/worker/slacksub"
	"github.com/lestrrat/go-pdebug"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
	"google.golang.org/cloud/pubsub"
)

type Router struct {
	*slacksub.Subscriber
	acmefwd chan *slack.MessageEvent
	deployfwd chan *slack.MessageEvent
}

func New(cl *pubsub.Client, topic, slackgwURL, token string) *Router {
	router := &Router{
		Subscriber: slacksub.New(cl, topic, slackgwURL, token),
		acmefwd:    make(chan *slack.MessageEvent, 128),
		deployfwd:  make(chan *slack.MessageEvent, 128),
	}
	router.Subscriber.MessageCallback = router.processMessageEvent

	go router.loop("acmebot-queue", router.acmefwd)
	go router.loop("deploybot-queue", router.deployfwd)

	return router
}

var spacesRx = regexp.MustCompile(`\s+`)

func (r *Router) processMessageEvent(ev *slack.MessageEvent) error {
	cmd := spacesRx.Split(strings.TrimSpace(ev.Text), -1)

	switch cmd[1] {
	case "acme":
		r.acmefwd <- ev
	case "deploy", "ingress":
		r.deployfwd <- ev
	default:
		if pdebug.Enabled {
			pdebug.Printf("Ignoring command '%s'", cmd[1])
		}
	}

	return nil
}

func (r *Router) loop(topicName string, inCh chan *slack.MessageEvent) {
	if pdebug.Enabled {
		g := pdebug.Marker("Start Router.loop(%s)", topicName)
		defer g.End()
	}

	flusht := time.Tick(time.Second)
	topic := r.Client.Topic(topicName)
	buf := make([]*slack.MessageEvent, 0, pubsub.MaxPublishBatchSize)
	msgs := make([]*pubsub.Message, 0, pubsub.MaxPublishBatchSize)
	for {
		select {
		case ev := <-inCh:
			buf = append(buf, ev)
			if len(buf) <= pubsub.MaxPublishBatchSize {
				continue
			}
		case <-flusht:
			if len(buf) == 0 {
				continue
			}
		}

		if pdebug.Enabled {
			pdebug.Printf("Processing %d events...", len(buf))
		}

		jsbuf := bytes.Buffer{}
		enc := json.NewEncoder(&jsbuf)
		for _, ev := range buf {
			jsbuf.Reset()
			if err := enc.Encode(slack.RTMEvent{Data: ev}); err != nil {
				if pdebug.Enabled {
					pdebug.Printf("ERROR: %s", err)
				}
				// Ugh. Ignore
				continue
			}
			msgs = append(msgs, &pubsub.Message{Data: jsbuf.Bytes()})
		}
		buf = buf[:0]

		// TODO: handle errors
		if pdebug.Enabled {
			pdebug.Printf("Forwarding %d messages to %s", len(msgs), topic.Name())
		}

		res, err := topic.Publish(context.Background(), msgs...)
		if pdebug.Enabled {
			if err != nil {
				pdebug.Printf("%s", err)
			}
			if res != nil {
				pdebug.Printf("%#v", res)
			}
		}
		msgs = msgs[:0]
	}
}