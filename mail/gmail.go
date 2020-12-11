package mail

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"text/template"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Mail structure
type Mail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	MIME    string `json:"mime"`
	Body    string `json:"body"`
}

// GmailService is the Gmail client for sending email
var GmailService *gmail.Service

// OAuthGmailService initializes the Gmail service
func OAuthGmailService(clientID, clientSecret, accessToken, refreshToken string) error {
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		return fmt.Errorf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		return fmt.Errorf("Email service is initialized")
	}

	return err
}

// SendEmailOAUTH2 sends a mail using Google's mail service
func SendEmailOAUTH2(to, subject string, data interface{}, template string) error {
	emailBody, err := parseTemplate(template, data)
	if err != nil {
		return fmt.Errorf("unable to parse email template - " + err.Error())
	}

	var message gmail.Message

	payload := Mail{
		To:      "To: " + to + "\r\n",
		Subject: "Subject: " + subject + "\n",
		MIME:    "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
		Body:    emailBody,
	}

	msg := []byte(payload.To + payload.Subject + payload.MIME + "\n" + payload.Body)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	_, err = GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return err
	}

	return err
}

// parseTemplate parses the gohtml template from the template folder
func parseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("mail/templates/%s", templateFileName))
	if err != nil {
		return "", fmt.Errorf("invalid template name")
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}
