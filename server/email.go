package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/smtp"
	"text/template"

	"github.com/ameske/go_nfl/database"
	"gopkg.in/yaml.v2"
)

var (
	config = emailConfig{}
	auth   smtp.Auth
	email  *template.Template
)

type emailConfig struct {
	EmailAddress    string `yaml:"EMAIL_ADDRESS"`
	Password        string `yaml:"PASSWORD"`
	SMTPAddress     string `yaml:"SMTP_ADDRESS"`
	SMTPFullAddress string `yaml:"SMTP_FULL_ADDRESS"`
}

type picksEmail struct {
	To      string
	From    string
	Subject string
	Week    int
	Picks   []database.FormPick
}

func LoadEmailConfig(path string) error {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return err
	}

	auth = smtp.PlainAuth("",
		config.EmailAddress,
		config.Password,
		config.SMTPAddress,
	)

	email = template.Must(template.ParseFiles("/opt/ameske/gonfl/templates/email.tmpl"))

	return err
}

func SendPicksEMail(to, subject string, week int, picks []database.FormPick) {
	pe := picksEmail{
		To:      to,
		From:    config.EmailAddress,
		Subject: subject,
		Week:    week,
		Picks:   picks,
	}

	var body bytes.Buffer
	email.Execute(&body, pe)

	to_s := make([]string, 0)
	to_s = append(to_s, config.EmailAddress)
	if config.EmailAddress != to {
		to_s = append(to_s, to)
	}

	err := smtp.SendMail(config.SMTPFullAddress, auth, config.EmailAddress, to_s, body.Bytes())
	if err != nil {
		log.Printf("Email Error: %s", err.Error())
	}
}
