package service

import (
	"os"
	"sync"

	"github.com/pkg/errors"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

var mailgunSvc *MailgunSvc
var mailgunOnce sync.Once

func Mailgun() *MailgunSvc {
	mailgunOnce.Do(mailgunSvc.Init)
	return mailgunSvc
}

func (v *MailgunSvc) Init() {
	f := func(v *string, envname string) {
		envvar := os.Getenv(envname)
		if envvar == "" {
			panic("Missing required environment variable " + envname)
		}
		*v = envvar
	}

	f(&v.defaultSender, "MAILGUN_DEFAULT_SENDER")

	var domain string
	var apiKey string
	var publicApiKey string
	f(&domain, "MAILGUN_DOMAIN")
	f(&apiKey, "MAILGUN_API_KEY")
	f(&publicApiKey, "MAILGUN_PUBLIC_API_KEY")

	v.client = mailgun.NewMailgun(domain, apiKey, publicApiKey)
}

type MailMessage struct {
	From       string
	Subject    string
	Text       string
	Recipients []string
}

func (v *MailgunSvc) Send(mm *MailMessage) error {
	if mm.From == "" {
		mm.From = v.defaultSender
	}

	m := mailgun.NewMessage(mm.From, mm.Subject, mm.Text, mm.Recipients...)

	mg := v.client
	_, _, err := mg.Send(m)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}
	return nil
}
