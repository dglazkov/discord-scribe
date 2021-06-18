package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/bwmarrin/discordgo"
	"github.com/dglazkov/discord-scribe/scribe"
)

const TOKEN_ENV_NAME = "BOT_TOKEN"
const APP_NAME = "discord-scribe"

func main() {
	token := os.Getenv(TOKEN_ENV_NAME)

	if token == "" {
		fmt.Println("No token provided. Please run: " + APP_NAME + " -t <bot token> or set env var " + TOKEN_ENV_NAME)
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	cfg := mysql.Cfg("flux-scribe-bot:us-west1:flux-scribe-database", "user", "")
	cfg.DBName = "flux"
	cfg.ParseTime = true

	const timeout = 10 * time.Second
	cfg.Timeout = timeout
	cfg.ReadTimeout = timeout
	cfg.WriteTimeout = timeout

	db, err := mysql.DialCfg(cfg)
	if err != nil {
		panic("couldn't dial: " + err.Error())
	}
	defer db.Close()

	// Register ready as a callback for the ready events.
	scribe := scribe.NewScribe(db, context.Background())
	newBot(dg, scribe)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}
	defer dg.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println(APP_NAME + " is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
