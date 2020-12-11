package main

import (
	"flag"
)

type config struct {
	APIHost string
	APIPort int

	DBUser string
	DBPass string
	DBName string
	DBHost string
	DBPort int

	HMAC string

	Salt string

	GmailClientID     string
	GmailClientSecret string
	GmailAccessToken  string
	GmailRefreshToken string
}

func newConfig() *config {
	c := &config{}

	flag.StringVar(&c.APIHost, "api_host", "127.0.0.1", "The API's host.")
	flag.IntVar(&c.APIPort, "api_port", 80, "The API's port.")

	flag.StringVar(&c.DBUser, "db_user", "root", "The database's username.")
	flag.StringVar(&c.DBPass, "db_pass", "root", "The database's password.")
	flag.StringVar(&c.DBName, "db_name", "database", "The database's name.")
	flag.StringVar(&c.DBHost, "db_host", "127.0.0.1", "The database's host.")
	flag.IntVar(&c.DBPort, "db_port", 3306, "The database's port.")

	flag.StringVar(&c.HMAC, "hmac", "", "HMAC secret")

	flag.StringVar(&c.Salt, "salt", "", "Password salt")

	flag.StringVar(&c.GmailClientID, "gmail_client_id", "", "The Google mail client ID")
	flag.StringVar(&c.GmailClientSecret, "gmail_client_secret", "", "The Google mail client secret")
	flag.StringVar(&c.GmailAccessToken, "gmail_access_token", "", "The Google mail access token")
	flag.StringVar(&c.GmailRefreshToken, "gmail_refresh_token", "", "The Google mail refresh token")

	flag.Parse()

	return c
}
