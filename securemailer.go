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

	mw := &MailWriter{}
	err = crypt.EncryptMessage(mw, from, to[0], string(body))
	if err != nil {
		log.Printf("Error encrypting message: %v\n", err)
		return
	}

	fmt.Printf("About to send:\n\nFrom: %s\nTo: %v\nBody: `%s`\n\n\n",
		from, to, mw.Enc)

	log.Println("TODO: Send email via Mailgun or the like")
}

func main() {
	lnAddr := "127.0.0.1:2525"
	log.Printf("Listening for emails on %s\n", lnAddr)
	log.Fatal(smtpd.ListenAndServe(lnAddr, mailHandler,
		"SecureMailer", ""))
}

type MailWriter struct {
	Enc []byte
}

func (mw *MailWriter) Write(p []byte) (int, error) {
	mw.Enc = append(mw.Enc, p...)
	return len(p), nil
}
