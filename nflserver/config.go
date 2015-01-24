package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Server ServerConfig
	Email  EmailConfig
}

type ServerConfig struct {
	AuthKey            string `json:"authKey"`
	EncryptKey         string `json:"encryptKey"`
	PostgresConnString string `json:"postgresConnString"`
}

type EmailConfig struct {
	SendAsAddress   string `json:"sendAsAddress"`
	Password        string `json:"password"`
	SMTPAddress     string `json:"smtpAddress"`
	SMTPFullAddress string `json:"smtpFullAddress"`
}

func loadConfig(path string) Config {
	configBytes, err := ioutil.ReadFile(path)

	config := Config{}
	err = json.Unmarshal(configBytes, &config)

	if err != nil {
		log.Fatalf(err.Error())
	}

	return config
}
