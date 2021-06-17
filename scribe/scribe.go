package scribe

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Scribe struct {
	db *sql.DB
}

func NewScribe(db *sql.DB) *Scribe {
	return &Scribe{db}
}

func asTime(s discordgo.Timestamp) time.Time {
	t, _ := time.Parse(time.RFC3339, string(s))
	return t
}

func (s *Scribe) AddMessage(message *discordgo.Message) {
	author_id := message.Author.ID
	_, err := s.db.Exec(`INSERT INTO
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
}
