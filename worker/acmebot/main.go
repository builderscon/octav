package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lestrrat/go-cloud-acmeagent"
	"github.com/lestrrat/go-cloud-acmeagent/gcp"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
	"google.golang.org/cloud/pubsub"
	"google.golang.org/cloud/storage"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	var authtokenf string
	var bucket string
	var email string
	var fifopath string
	var projectID string
	var slackgw string
	var topic string
	var zone string
	flag.StringVar(&authtokenf, "authtokenfile", "", "File containing token used to authentication when posting")
	flag.StringVar(&bucket, "bucket", "", "bucket name (default is 'acme-' + projectID)")
	flag.StringVar(&email, "email", "", "email ID to use for acme protocol")
	flag.StringVar(&fifopath, "fifopath", "", "path to where tls requests willbe pushed to")
	flag.StringVar(&projectID, "project_id", "", "project ID to use")
	flag.StringVar(&slackgw, "slackgw", "http://slackgw:4979", "slack gateway url")
	flag.StringVar(&topic, "topic", "slackgw-url", "topic name to subscribe to")
	flag.StringVar(&zone, "zone", "", "DNS zone to update")
	flag.Parse()

	if fifopath == "" {
		fmt.Printf("fifopath is required")
		return 1
	}

	if projectID == "" {
		fmt.Printf("projectID is required")
		return 1
	}

	if bucket == "" {
		bucket = "acme-" + projectID
	}

	var authtoken string
	if authtokenf != "" {
		buf, err := ioutil.ReadFile(authtokenf)
		if err != nil {
			fmt.Printf("failed to read file '%s': %s", authtokenf, err)
			return 1
		}
		authtoken = string(buf)
	}

	pcl, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		fmt.Printf("failed to create pubsub client: %s", err)
		return 1
	}

	scl, err := storage.NewClient(context.Background())
	if err != nil {
		fmt.Printf("failed to create storage client: %s", err)
		return 1
	}

	store := gcp.NewStorage(scl, projectID, email, bucket)

	httpcl, err := google.DefaultClient(
		context.Background(),
		dns.NdevClouddnsReadwriteScope, // We need to be able to update CloudDNS
	)
	if err != nil {
		panic(err)
	}

	dnssvc, err := dns.New(httpcl)
	if err != nil {
		panic(err)
	}

	aa, err := acmeagent.New(acmeagent.AgentOptions{
		DNSCompleter: gcp.NewDNS(dnssvc, projectID, zone),
		StateStorage: store,
	})

	bot := New(pcl, aa, store, topic, slackgw, authtoken, fifopath)
	bot.Run()

	return 0
}