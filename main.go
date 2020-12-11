package main

import (
	"cat-clerk-api/api"
	"cat-clerk-api/auth"
	"cat-clerk-api/database"
	"cat-clerk-api/mail"
	"cat-clerk-api/util"
	"fmt"
	"log"
	"net/http"
	"os"

	muxHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	cfg := newConfig()

	db := database.Init(cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBHost, cfg.DBPort)

	if err := db.EnsureConnected(); err != nil {
		log.Fatal(err)
		return
	}

	util.Init(cfg.Salt)

	mail.OAuthGmailService(
		cfg.GmailClientID,
		cfg.GmailClientSecret,
		cfg.GmailAccessToken,
		cfg.GmailRefreshToken,
	)

	router := mux.NewRouter().StrictSlash(true)

	auth := auth.New(router, []byte(cfg.HMAC))

	restAPI := api.Init(router, db, auth)

	router = restAPI.Handlers()

	requestLogger := muxHandlers.LoggingHandler(os.Stdout, auth)

	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%d", cfg.APIPort),
			requestLogger,
		),
	)
}
