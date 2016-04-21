## SecureMailer

SecureMailer is a secure SMTP relay server.  It accepts plaintext
emails over SMTP, encrypts them with the recipient's GPG key, signs
them with the sender's GPG key, then delivers them through Mailgun.


### Pre-requisites

Currently, the sender's key must not be passphrase-protected, and the
recpient's key must be in the local GPG keyring (at
`~/.gnupg/pubring.gpg`).


### Getting Started

```
go get github.com/elimisteve/securemailer
MAILGUN_DOMAIN="mg.mydomain.com" MAILGUN_API_KEY="key-..." securemailer
```

Optionally, set the `PORT` environment variable to determine the port
that SecureMailer listens on (defaults to 2525).


### Quality

Alpha quality, mostly due to the underlying GPG library I wrote years
ago.  Please [create an
issue](https://github.com/elimisteve/securemailer/issues) if you find
a bug.


### Why SecureMailer?

SecureMailer was created to send passwordless login emails to
[Sandstorm](https://sandstorm.io/) users who need higher levels of
privacy than unencrypted emails allow, but should be useful in various
situations.


### Why Mailgun?

Because it's cheap, seems to be reliable, and Rackspace (who owns
Mailgun) isn't a
[PRISM](https://en.wikipedia.org/wiki/PRISM_%28surveillance_program%29)
partner.  I am very open to pull requests that add support for other
mailer services.
