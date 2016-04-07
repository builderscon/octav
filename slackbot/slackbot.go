package slackbot

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/lestrrat/go-cloud-acmeagent"
	"github.com/lestrrat/go-cloud-acmeagent/gcp"
	"github.com/lestrrat/go-pdebug"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/storage/v1"
	k8sc "k8s.io/kubernetes/pkg/client/unversioned"
)

var slackClient *slack.Client
var slackUser string
var acmeAgent *acmeagent.AcmeAgent
var acmeStateStore acmeagent.StateStorage

func init() {
	if err := initSlack(); err != nil {
		panic(err)
	}

	if err := initAcmeAgent(); err != nil {
		panic(err)
	}
}

func initSlack() error {
	var token string
	if err := readEnvConfigFile("Slack API token", "SLACKBOT_API_TOKEN_FILE", &token); err != nil {
		return err
	}
	if token == "" {
		return errors.New("token is empty")
	}

	slackClient = slack.New(token)
	autht, err := slackClient.AuthTest()
	if err != nil {
		return err
	}
	pdebug.Printf("%#v", autht)
	slackUser = autht.UserID
	return nil
}

func initAcmeAgent() error {
	email := os.Getenv("ACME_AGENT_EMAIL")
	if email == "" {
		return errors.New("ACME_AGENT_EMAIL environment variable is required for this test")
	}
	gcpproj := os.Getenv("ACME_AGENT_GCP_PROJECT_ID")
	if gcpproj == "" {
		return errors.New("ACME_AGENT_GCP_PROJECT_ID environment variable is required for this test")
	}
	gcpzone := os.Getenv("ACME_AGENT_GCP_ZONE_NAME")
	if gcpzone == "" {
		return errors.New("ACME_AGENT_GCP_ZONE_NAME environment variable is required for this test")
	}

	ctx := context.Background()
	httpcl, err := google.DefaultClient(ctx,
		dns.NdevClouddnsReadwriteScope,
		storage.CloudPlatformScope,
		storage.DevstorageReadWriteScope,
	)
	if err != nil {
		return err
	}
	storagesvc, err := storage.New(httpcl)
	if err != nil {
		return err
	}
	acmeStateStore = gcp.NewStorage(storagesvc, gcpproj, email, "acme-"+gcpproj)

	dnssvc, err := dns.New(httpcl)
	if err != nil {
		return err
	}

	k8sClient, err := k8sc.NewInCluster()
	if err != nil {
		return err
	}

	aa, err := acmeagent.New(acmeagent.AgentOptions{
		DNSCompleter: gcp.NewDNS(dnssvc, gcpproj, gcpzone),
		Uploader:     gcp.NewSecretUpload(k8sClient, "default"),
		StateStorage: acmeStateStore,
	})

	acmeAgent = aa
	return nil
}

// Dummy for now
func Run(_ string) error {
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(2)

	go StartRTM(done, &wg)
	go StartWatch(done, &wg)

	sigCh := make(chan os.Signal, 265)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	for {
		select {
		case <-sigCh:
			if pdebug.Enabled {
				pdebug.Printf("Received signal, bailing out")
			}
			close(done)
		}
	}

	if pdebug.Enabled {
		pdebug.Printf("Waiting for goroutines to exit...")
	}
	wg.Wait()

	if pdebug.Enabled {
		pdebug.Printf("All goroutines have exited, bailing out")
	}
	return nil
}
