package scholaemail

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	FromAddress         = "admin@scholacantorum.org"
	FromAddressWithName = "Schola Cantorum <admin@scholacantorum.org>"
)

// Send sends an email through the admin@scholacantorum.org account.  The
// supplied buffer must be a complete, RFC-2822-encoded email message, with
// the intended recipients expressed in To:, Cc:, and Bcc: headers.  It does
// not have to have a Date: header; one will be automatically added.
//
// For this to work, the appropriate Google API credentials must exist in
// $HOME/.config/gcloud/application_default_credentials.json.  Since $HOME
// might not be set, there is a fallback to the first two components of $PWD
// (i.e., usually /home/XXX or /Users/XXX).
func Send(data []byte) (err error) {
	var (
		key  []byte
		conf *jwt.Config
		ts   oauth2.TokenSource
		svc  *gmail.Service
		msg  *gmail.Message
	)
	if key, err = os.ReadFile(home() + "/.config/gcloud/application_default_credentials.json"); err != nil {
		return fmt.Errorf("sendmail: can't read Google API key: %w", err)
	}
	if conf, err = google.JWTConfigFromJSON(key, gmail.GmailSendScope); err != nil {
		return fmt.Errorf("sendmail: can't use Google API key: %w", err)
	}
	conf.Subject = FromAddress
	ts = conf.TokenSource(context.Background())
	if svc, err = gmail.NewService(context.Background(), option.WithTokenSource(ts)); err != nil {
		return fmt.Errorf("sendmail: can't create Gmail service client: %w", err)
	}
	msg = &gmail.Message{Raw: base64.RawURLEncoding.EncodeToString(data)}
	if _, err = svc.Users.Messages.Send(FromAddress, msg).Do(); err != nil {
		return fmt.Errorf("sendmail: %w", err)
	}
	return nil
}

func home() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	if cwd, err := os.Getwd(); err == nil {
		if parts := strings.Split(cwd, string(filepath.Separator)); parts[0] == "" && len(parts) >= 3 {
			return strings.Join(parts[0:3], string(filepath.Separator))
		}
	}
	return "" // we won't find the JSON file; the mail send will fail
}
