package main

import (
	"bytes"
	"log"
	"net/smtp"
	"text/template"

	"github.com/ameske/nfl-pickem/database"
)

var (
	auth       smtp.Auth
	sendAddr   string
	smtpServer string
	email      *template.Template
)

type picksEmail struct {
	To      string
	From    string
	Subject string
	Week    int
	Picks   []database.FormPick
}

func configureEmail(config Config) {
	auth = smtp.PlainAuth("",
		config.Email.SendAsAddress,
		config.Email.Password,
		config.Email.SMTPAddress,
	)

	email = template.Must(template.ParseFiles("/opt/ameske/gonfl/templates/email.tmpl"))

	sendAddr = config.Email.SendAsAddress
	smtpServer = config.Email.SMTPFullAddress
}

func SendPicksEmail(to, subject string, week int, picks []database.FormPick) {
	pe := picksEmail{
		To:      to,
		From:    sendAddr,
		Subject: subject,
		Week:    week,
		Picks:   picks,
	}

	var body bytes.Buffer
	email.Execute(&body, pe)

	to_s := make([]string, 0)
	to_s = append(to_s, sendAddr)
	if sendAddr != to {
		to_s = append(to_s, to)
	}

	err := smtp.SendMail(smtpServer, auth, sendAddr, to_s, body.Bytes())
	if err != nil {
		log.Printf("Email Error: %s", err.Error())
	}
}
