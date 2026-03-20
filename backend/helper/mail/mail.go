// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package mail

import (
	"bytes"
	"embed"
	"encoding/base64"
	"strings"

	//"html/template"
	"net/smtp"
	"text/template"

	"github.com/eclipse-disuko/disuko/conf"
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

var templates []string = []string{"taskReview", "taskApproval"}

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

func (c Client) Send(to string, templateName string, data any) error {
	if conf.Config.Server.E2ETests {
		return nil
	}
	if c.Host == "" {
		return nil
	}

	tmpl, err := template.New("taskEmail").ParseFS(taskTemplate, "templates/"+templateName+".txt")
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	/*
		htmlBody := new(bytes.Buffer)
		err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
		if err != nil {
			return err
		}
	*/

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + explicitUTF8Subject(subject.String()) + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"Content-Transfer-Encoding: 8bit\r\n" +
		"\r\n" +
		plainBody.String() + "\r\n")

	var auth smtp.Auth
	if c.User != "" && c.Pass != "" {
		auth = smtp.PlainAuth("", c.User, c.Pass, c.Host)
	}
	err = smtp.SendMail(c.Host+":"+c.Port, auth, c.Sender, []string{to}, msg)
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
