package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/cloud/pubsub"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	var authtokenf string
	var projectID string
	var slackgw string
	var topic string
	flag.StringVar(&authtokenf, "authtokenfile", "", "File containing token used to authentication when posting")
	flag.StringVar(&projectID, "project_id", "", "project ID to use")
	flag.StringVar(&slackgw, "slackgw", "http://slackgw:4979", "slack gateway url")
	flag.StringVar(&topic, "topic", "slackgw-forward", "topic name to subscribe to")
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

	pcl, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		fmt.Printf("failed to create pubsub client: %s", err)
		return 1
	}

	bot := New(pcl, topic, slackgw, authtoken)
	bot.Run()

	return 0
}