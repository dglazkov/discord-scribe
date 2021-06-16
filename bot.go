package main

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	"github.com/bwmarrin/discordgo"
)

type bot struct {
	session *discordgo.Session
	db      *sql.DB
}

func newBot(s *discordgo.Session, db *sql.DB) *bot {
	result := &bot{
		session: s,
		db:      db,
	}
	s.AddHandler(result.ready)
	s.AddHandler(result.messageCreate)
	return result
}

// discordgo callback: called when the bot receives the "ready" event from Discord.
func (b *bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	//GuildInfo isn't populated yet.
	fmt.Println("Ready and waiting!")
}

func asTime(s discordgo.Timestamp) time.Time {
	t, _ := time.Parse(time.RFC3339, string(s))
	return t
}

// discordgo callback: called after the when new message is posted.
func (b *bot) messageCreate(s *discordgo.Session, event *discordgo.MessageCreate) {
	message := event.Message
	author_id := message.Author.ID
	_, err := b.db.Exec(`INSERT INTO
		messages (id, channel_id, guild_id, author_id, content, timestamp)
		values (?, ?, ?, ?, ?, ?)`,
		message.ID,
		message.ChannelID,
		message.GuildID,
		author_id,
		message.Content,
		asTime(message.Timestamp))

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Message created")
	// var now time.Time
	// fmt.Println(b.db.QueryRow("SELECT NOW()").Scan(&now))
	// fmt.Println(now)
}
