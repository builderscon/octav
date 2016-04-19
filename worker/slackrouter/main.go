package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/cloud/pubsub"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	var projectID string
	var slackgw string
	var topic string
	flag.StringVar(&projectID, "project_id", "", "project ID to use")
	flag.StringVar(&slackgw, "slackgw", "http://slackgw:4979", "slack gateway url")
	flag.StringVar(&topic, "topic", "slackgw-forward", "topic name to subscribe to")
	flag.Parse()

	pcl, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		fmt.Printf("failed to create pubsub client: %s", err)
		return 1
	}

	bot := New(pcl, topic, slackgw)
	bot.Run()

	return 0
}