# Slackbot

## Notify GKE Pod status

Automatic detection of pods going up or down are notified to Slack

## Let's Encrypt Support

Simply talk to the bot to create and upload certificates for builderscon.io pages.

### `acme authz <domain>`

![](../../media/images/slackbot-letsencrypt-authz.png)

Get authorization to issue certificates for given domain. We use dns-01
challenge, so this takes approx 5 to 10 minutes due to the way CloudDNS
works.

### `acme cert <domain>`

Fetch certificates for the given domain.

### `acme upload <domain>`

![](../../media/images/slackbot-letsencrypt-upload.png)

Upload new certificates that then can be used in HTTPS Load Balancers.