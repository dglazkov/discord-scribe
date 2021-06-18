package scribe

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Scribe struct {
	db  *sql.DB
	ctx context.Context
}

func NewScribe(db *sql.DB, ctx context.Context) *Scribe {
	return &Scribe{db, ctx}
}

func asTime(s discordgo.Timestamp) time.Time {
	t, _ := time.Parse(time.RFC3339, string(s))
	return t
}

func (s *Scribe) AddMessage(message *discordgo.Message) {
	author_id := message.Author.ID
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}

	_, exec_err := tx.Exec(`INSERT INTO
		messages (id, channel_id, guild_id, author_id, content, timestamp)
		values (?, ?, ?, ?, ?, ?)`,
		message.ID,
		message.ChannelID,
		message.GuildID,
		author_id,
		message.Content,
		asTime(message.Timestamp))

	if exec_err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message added")
}
