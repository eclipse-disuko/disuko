// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package mail

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"strings"

	//"html/template"
	"net/smtp"
	"text/template"
)

//go:embed templates
var taskTemplate embed.FS

type Client struct {
	Host   string
	Port   string
	Sender string
	User   string
	Pass   string
}

var templates []string = []string{"taskReview", "taskApproval", "approvalInactiveMail"}

func NewClient(host, port, sender, user, pass string) Client {
	return Client{
		Host:   host,
		Port:   port,
		Sender: sender,
		User:   user,
		Pass:   pass,
	}
}

// explicitUTF8Subject encodes the email subject using RFC 2047 MIME encoding.
// This ensures proper handling of UTF-8 characters (like umlauts, accents, etc.)
// in email subject lines by base64-encoding them with the proper charset declaration.
// Format: =?UTF-8?B?<base64-encoded-text>?=
func explicitUTF8Subject(subject string) string {
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?="
}

func (c Client) Send(to string, cc string, bcc string, bodyTmpl, subjectTmpl string, data any) error {
	if c.Host == "" {
		return nil
	}

	t, err := template.New("subject").Parse(subjectTmpl)
	if err != nil {
		return fmt.Errorf("parsing subject template: %w", err)
	}
	subject := new(bytes.Buffer)
	err = t.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return fmt.Errorf("executing subject template: %w", err)
	}

	t, err = template.New("body").Parse(bodyTmpl)
	if err != nil {
		return fmt.Errorf("parsing body template: %w", err)
	}
	body := new(bytes.Buffer)
	err = t.ExecuteTemplate(body, "body", data)
	if err != nil {
		return fmt.Errorf("executing body template: %w", err)
	}

	msgLines := []string{
		"To: " + to,
	}
	if cc != "" {
		msgLines = append(msgLines, "CC: "+cc)
	}
	msgLines = append(
		msgLines,
		"Subject: "+explicitUTF8Subject(subject.String()),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Transfer-Encoding: 8bit",
		"",
		body.String(),
	)

	var auth smtp.Auth
	if c.User != "" && c.Pass != "" {
		auth = smtp.PlainAuth("", c.User, c.Pass, c.Host)
	}

	toAddrs := []string{to}
	if bcc != "" {
		toAddrs = append(toAddrs, bcc)
	}

	err = smtp.SendMail(c.Host+":"+c.Port, auth, c.Sender, toAddrs, []byte(strings.Join(msgLines, "\r\n")))
	return err
}

func (c Client) IsTeamplateValid(templateName string) bool {
	if templateName == "" {
		return false
	}
	for _, template := range templates {
		if strings.EqualFold(template, templateName) {
			return true
		}
	}
	return false
}
