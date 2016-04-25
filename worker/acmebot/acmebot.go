package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/builderscon/octav/worker/acmebot/slacksub"
	"github.com/lestrrat/go-cloud-acmeagent"
	"github.com/lestrrat/go-jwx/jwk"
	"github.com/lestrrat/go-pdebug"
	"github.com/nlopes/slack"
	"google.golang.org/cloud/pubsub"
)

type Bot struct {
	*slacksub.Subscriber
	acmeagent *acmeagent.AcmeAgent
	acmestore acmeagent.StateStorage
	fifopath  string
}

type acmectx struct {
	msg *slack.MessageEvent
}

func New(cl *pubsub.Client, agent *acmeagent.AcmeAgent, store acmeagent.StateStorage, topic, slackgwURL, authtoken, fifopath string) *Bot {
	bot := &Bot{
		Subscriber: slacksub.New(cl, topic, slackgwURL, authtoken),
		acmeagent:  agent,
		acmestore:  store,
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

	if cmd[1] != "acme" {
		return nil
	}
	ctx := acmectx{
		msg: ev,
	}
	return b.handleLetsEncryptCmd(&ctx, cmd[2:])
}

func (b *Bot) handleLetsEncryptCmd(ctx *acmectx, cmd []string) error {
	switch cmd[0] {
	case "help":
		return b.handleHelpCmd(ctx)
	case "authz":
		if len(cmd) < 2 {
			return b.handleHelpCmd(ctx)
		}
		return b.handleAuthzCmd(ctx, cmd[1:])
	case "cert":
		if len(cmd) < 3 {
			return b.handleHelpCmd(ctx)
		}
		return b.handleCertCmd(ctx, cmd[1:])
	}
	return b.handleHelpCmd(ctx)
}

func (b *Bot) reply(ctx *acmectx, message string) error {
	return b.postMessage(ctx, ctx.msg.Channel, message)
}

func (b *Bot) postMessage(ctx *acmectx, channel, message string) error {
	_, err := http.PostForm(
		b.SlackgwURL+"/post",
		url.Values{
			"channel": []string{channel},
			"message": []string{message},
		},
	)
	return err
}

func (b *Bot) handleHelpCmd(ctx *acmectx) error {
	return b.reply(ctx, `usage: acme [cert|authz] [subcmds...]

acme cert issue <domain>
acme cert delete <domain>
acme cert upload <domain>
acme authz request <domain>
acme authz delete <domain>
`)
}

func (b *Bot) handleAuthzCmd(ctx *acmectx, cmd []string) error {
	if len(cmd) < 2 {
		return b.handleHelpCmd(ctx)
	}

	sl, err := parseSlackLink(cmd[1])
	if err != nil {
		return err
	}
	domain := sl.Text
	switch cmd[0] {
	case "request":
		return b.handleAuthzRequestCmd(ctx, domain)
	case "delete":
		return b.handleAuthzDeleteCmd(ctx, domain)
	case "show":
		return b.handleAuthzShowCmd(ctx, domain)
	default:
		return b.handleHelpCmd(ctx)
	}
}

func (b *Bot) handleAuthzDeleteCmd(ctx *acmectx, domain string) error {
	if err := b.acmestore.DeleteAuthorization(domain); err != nil {
		return b.reply(ctx, ":exclamation: Deleting authorization failed: "+err.Error())
	}
	return b.reply(ctx, ":tada: Deleted authorization")
}

func (b *Bot) handleAuthzRequestCmd(ctx *acmectx, domain string) error {
	b.reply(ctx, ":white_check_mark: Authorizing *"+domain+"*")

	var authz acmeagent.Authorization
	if err := b.acmestore.LoadAuthorization(domain, &authz); err != nil {
		b.reply(ctx, ":white_check_mark: Authorization for domain not found in storage.")
	} else {
		if authz.IsExpired() {
			b.reply(ctx, ":exclamation: Authorization expired, going to run authorization again")
		} else {
			return b.reply(ctx, ":exclamation: Authorization already exists. Run `acme cert` to issue certificates for this domain")
		}
	}

	b.reply(ctx, ":white_check_mark: Running authorization (this may take a few minutes)")

	if err := b.acmeagent.AuthorizeForDomain(domain); err != nil {
		return b.reply(ctx, ":exclamation: Authorization failed: "+err.Error())
	}
	return b.reply(ctx, ":tada: Authorization for domain *"+domain+"* complete")
}

func (b *Bot) handleAuthzShowCmd(ctx *acmectx, domain string) error {
	var authz acmeagent.Authorization
	if err := b.acmestore.LoadAuthorization(domain, &authz); err != nil {
		return b.reply(ctx, ":white_check_mark: Authorization for domain not found in storage.")
	}

	buf, _ := json.MarshalIndent(authz, "", "  ")
	return b.reply(ctx, "```\n"+string(buf)+"\n```")
}

func (b *Bot) handleCertCmd(ctx *acmectx, cmd []string) error {
	switch len(cmd) {
	case 0, 1:
		return b.reply(ctx, "Usage: `acme cert [issue|delete|upload] <domain>`")
	default:
	}

	sl, err := parseSlackLink(cmd[1])
	if err != nil {
		return err
	}
	domain := sl.Text

	switch cmd[0] {
	case "issue":
		return b.handleCertIssueCmd(ctx, domain)
	case "delete":
		return b.handleCertDeleteCmd(ctx, domain)
	case "upload":
		return b.handleCertUploadCmd(ctx, domain)
	default:
		return b.reply(ctx, "Usage: `acme cert [issue|delete|upload] <domain>`")
	}
}

func (b *Bot) handleCertDeleteCmd(ctx *acmectx, domain string) error {
	b.reply(ctx, ":white_check_mark: Deleting certificates for *"+domain+"*")
	if err := b.acmestore.DeleteCert(domain); err != nil {
		return b.reply(ctx, ":exclamation: Failed to delete certificates: "+err.Error())
	}
	return b.reply(ctx, ":tada: Deleted certificates")
}

func (b *Bot) handleCertIssueCmd(ctx *acmectx, domain string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("b.handleCertIssue %s", domain).BindError(&err)
		defer g.End()
	}

	b.reply(ctx, ":white_check_mark: Issueing certificates for *"+domain+"*")

	var cert *x509.Certificate
	if err := b.acmestore.LoadCert(domain, cert); err != nil {
		b.reply(ctx, ":white_check_mark: Certificates for domain not found in storage.")
	} else {
		if time.Now().After(cert.NotAfter) {
			b.reply(ctx, ":exclamation: Certificate expired, going to issue it again")
		} else {
			return b.reply(ctx, ":exclamation: Certificate already exists. Run `acme upload` to upload the certificate")
		}
	}

	// run handleAuthzCmd to make sure that the authorization is there
	if err := b.handleAuthzRequestCmd(ctx, domain); err != nil {
		return err
	}

	b.reply(ctx, ":white_check_mark: Fetching certificates")
	if err := b.acmeagent.IssueCertificate(domain, nil, false); err != nil {
		return b.reply(ctx, ":exclamation: Failed to fetch certificates: "+err.Error())
	}
	return b.reply(ctx, ":tada: Issueing certificates for domain *"+domain+"* complete")
}

func (b *Bot) handleCertUploadCmd(ctx *acmectx, domain string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("b.handleCertUploadCmd %s", domain).BindError(&err)
		defer g.End()
	}

	b.reply(ctx, ":white_check_mark: Uploading certificates for *"+domain+"*")

	// Instead of uploading somewhere, write to a fifo so another
	// process can read and process k8s stuff. This is done to avoid
	// linking programs that use google.golang.org/cloud with libraries
	// from k8s.io/kubernetes

	// Load the cert and key
	var cert x509.Certificate
	var key jwk.RsaPrivateKey
	if err := b.acmestore.LoadCertFullChain(domain, &cert); err != nil {
		return err
	}
	if err := b.acmestore.LoadCertKey(domain, &key); err != nil {
		return err
	}

	// pem encode them certificates and keys
	var pemcert bytes.Buffer
	var pemkey bytes.Buffer
	if err := pem.Encode(&pemcert, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}); err != nil {
		return err
	}

	privkey, err := key.PrivateKey()
	if err != nil {
		return err
	}

	if err := pem.Encode(&pemkey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privkey)}); err != nil {
		return err
	}

	fifo, err := b.Fifo()
	if err != nil {
		return err
	}
	defer fifo.Close()

	// Write to local fifo
	name := domain + "-" + time.Now().Format("20060102-150405")
	err = json.NewEncoder(fifo).Encode(map[string]string{
		"name":    name,
		"channel": ctx.msg.Channel,
		"tls.crt": pemcert.String(),
		"tls.key": pemkey.String(),
	})
	if err != nil {
		return err
	}

	return b.reply(ctx, ":white_check_mark: Certificate has been sent to be processed. Please wait a moment for the secret '"+name+"' to be available")
}