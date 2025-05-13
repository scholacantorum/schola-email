// send-raw-email sends an email from the admin@scholacantorum.org account.
// The headers and body of the message are given through standard input.
package main

import (
	"bytes"
	"io"
	"os"

	belog "github.com/scholacantorum/public-site-backend/backend-log"
	scholaemail "github.com/scholacantorum/schola-email"
)

func main() {
	var (
		buf bytes.Buffer
		err error
	)
	belog.LogApp = "send-raw-email"
	if _, err = io.Copy(&buf, os.Stdin); err != nil {
		belog.Log("can't read input: %s", err)
		os.Exit(1)
	}
	if err = scholaemail.Send(buf.Bytes()); err != nil {
		belog.Log("%s", err)
		os.Exit(1)
	}
}
