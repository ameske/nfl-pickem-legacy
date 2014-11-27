package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	emailAddress    string `yaml:"EMAIL_ADDRESS"`
	password        string `yaml:"PASSWORD"`
	smtpAddress     string `yaml:"SMTP_ADDRESS"`
	smtpFullAddress string `yaml:"SMTP_FULL_ADDRESS"`
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
		config.emailAddress,
		config.password,
		config.smtpAddress,
	)

	email = template.Must(template.ParseFiles("email.tmpl"))

	return err
}

func SendPicksEMail(to, subject string, week int, picks []database.FormPick) {
	pe := picksEmail{
		To:      to,
		From:    config.emailAddress,
		Subject: subject,
		Week:    week,
		Picks:   picks,
	}

	var body bytes.Buffer
	email.Execute(&body, pe)

	ioutil.WriteFile(fmt.Sprintf("%s - %s.txt", to, subject), body.Bytes(), 777)
	//smtp.SendMail(config.smtpFullAddress, auth, config.emailAddress, []string{to}, body.Bytes())
}
