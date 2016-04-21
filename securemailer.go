// Steven Phillips / elimisteve
// 2016.04.20

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"os"

	mailgun "github.com/mailgun/mailgun-go"
	"github.com/mhale/smtpd"
	"github.com/thecloakproject/utils/crypt"
)

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	if len(to) == 0 {
		log.Println("FATAL ERROR: email must be sent to someone; returning early")
		return
	}

	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Printf("FATAL ERROR parsing email; returning early: %v\n", err)
		return
	}

	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Printf("FATAL ERROR reading message body; returning early: %v\n", err)
		return
	}

	// TODO: If safe to do so, prepend "Subject: ...\n\n" to
	// original message body before encryption

	encMsg, err := encryptMessage(from, to[0], string(body))
	if err != nil {
		log.Printf("Error encrypting message: %v\n", err)
		return
	}

	fmt.Printf("About to send:\n\nFrom: %s\nTo: %v\nBody: `%s`\n\n\n",
		from, to, encMsg)

	resp, id, err := mailgunSend(from, to[0], string(encMsg))
	if err != nil {
		log.Printf("Error sending via Mailgun: %v\n", err)
		return
	}

	log.Printf("Mailgun id: %v; Response: %v\n", id, resp)
}

func main() {
	lnAddr := "127.0.0.1:2525"
	log.Printf("Listening for emails on %s\n", lnAddr)
	log.Fatal(smtpd.ListenAndServe(lnAddr, mailHandler,
		"SecureMailer", ""))
}

var gun = mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"),
	os.Getenv("MAILGUN_API_KEY"), os.Getenv("MAILGUN_PUBLIC_KEY"))

func mailgunSend(from, to, body string) (resp, id string, err error) {
	m := mailgun.NewMessage(from, "", body, to)
	return gun.Send(m)
}

type MailWriter struct {
	Enc []byte
}

func (mw *MailWriter) Write(p []byte) (int, error) {
	mw.Enc = append(mw.Enc, p...)
	return len(p), nil
}

func encryptMessage(from, to, body string) (enc []byte, err error) {
	mw := &MailWriter{}
	err = crypt.EncryptMessage(mw, from, to, body)
	if err != nil {
		return nil, err
	}
	return mw.Enc, nil
}
