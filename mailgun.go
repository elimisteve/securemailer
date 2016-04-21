// Steven Phillips / elimisteve
// 2016.04.21

package main

import (
	"log"
	"os"

	mailgun "github.com/mailgun/mailgun-go"
)

var gun = mailgun.NewMailgun(
	os.Getenv("MAILGUN_DOMAIN"),
	os.Getenv("MAILGUN_API_KEY"),
	os.Getenv("MAILGUN_PUBLIC_KEY"),
)

func init() {
	if gun.Domain() == "" {
		log.Fatalln("Mailgun misconfigured: no domain set. Set MAILGUN_DOMAIN.")
	}
	if gun.ApiKey() == "" {
		log.Fatalln("Mailgun misconfigured: no API key. Set MAILGUN_API_KEY.")
	}
}

func mailgunSend(from, to, subject, body string) (resp, id string, err error) {
	m := mailgun.NewMessage(from, subject, body, to)
	return gun.Send(m)
}
