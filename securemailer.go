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

	"github.com/mhale/smtpd"
)

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	if len(to) == 0 {
		log.Println("Error: email must be sent to someone; returning early")
		return
	}

	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Printf("Error parsing email; returning early: %v\n", err)
		return
	}

	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Printf("Error reading message body; returning early: %v\n", err)
		return
	}

	// TODO: If safe to do so, prepend "Subject: ...\n\n" to
	// original message body before encryption

	encMsg, err := encryptMessage(from, to[0], string(body))
	if err != nil {
		log.Printf("Error encrypting message: %v\n", err)
		return
	}

	fmt.Printf("About to send:\n\nFrom: %s\nTo: %v\nBody: `%s`\n\n",
		from, to[0], encMsg)

	// TODO: Make it configurable whether the original subject is
	// included
	subject := ""

	resp, id, err := mailgunSend(from, to[0], subject, string(encMsg))
	if err != nil {
		log.Printf("Error sending via Mailgun: %v\n", err)
		return
	}

	log.Printf("Mailgun id: %v; Response: %v\n\n\n", id, resp)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "2525"
	}
	listenAddr := "127.0.0.1:" + port

	log.Printf("Listening for emails on %s\n", listenAddr)
	log.Fatal(smtpd.ListenAndServe(listenAddr, mailHandler,
		"SecureMailer", ""))
}
