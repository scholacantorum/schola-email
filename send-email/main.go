// send-email sends an email from the admin@scholacantorum.org account.
// The headers and body of the message are given through standard input.
// The body is interpreted as HTML, and is wrapped in Schola styles and logo.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"os"

	belog "github.com/scholacantorum/public-site-backend/backend-log"
	scholaemail "github.com/scholacantorum/schola-email"
)

func main() {
	var buf bytes.Buffer
	var mw *multipart.Writer
	var html io.Writer
	var htmlqp io.Writer
	var img io.Writer
	var input *bufio.Scanner
	var hdr textproto.MIMEHeader
	var err error

	belog.LogApp = "send-email"
	input = bufio.NewScanner(os.Stdin)
	for input.Scan() {
		if input.Text() == "" {
			break
		}
		buf.WriteString(input.Text())
		buf.WriteString("\r\n")
	}
	mw = multipart.NewWriter(&buf)
	fmt.Fprintf(&buf, "Content-Type: multipart/related; boundary=%s\r\n\r\n", mw.Boundary())
	hdr = make(textproto.MIMEHeader)
	hdr.Set("Content-Type", "text/html; charset=UTF-8")
	hdr.Set("Content-Transfer-Encoding", "quoted-printable")
	html, _ = mw.CreatePart(hdr)
	htmlqp = quotedprintable.NewWriter(html)
	htmlqp.Write([]byte(`<!DOCTYPE html><html><body style="margin:0"><div style="width:600px;margin:0 auto"><div style="margin-bottom:24px"><img src="cid:SCHOLA_LOGO" alt="[Schola Cantorum]" style="border-width:0"></div>`))
	for input.Scan() {
		htmlqp.Write(input.Bytes())
		htmlqp.Write([]byte{'\n'})
	}
	if err = input.Err(); err != nil {
		belog.Log("can't read input: %s", err)
		os.Exit(1)
	}
	htmlqp.Write([]byte("</div></body></html>\n"))
	hdr = make(textproto.MIMEHeader)
	hdr.Set("Content-Type", "image/gif")
	hdr.Set("Content-Transfer-Encoding", "base64")
	hdr.Set("Content-ID", "<SCHOLA_LOGO>")
	img, _ = mw.CreatePart(hdr)
	img.Write([]byte(mailLogo))
	mw.Close()
	if err = scholaemail.Send(buf.Bytes()); err != nil {
		belog.Log("%s", err)
		os.Exit(1)
	}
}
