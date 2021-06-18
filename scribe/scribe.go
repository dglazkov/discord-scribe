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
	db     *sql.DB
	ctx    context.Context
	reader MessageReader
}

func NewScribe(db *sql.DB, ctx context.Context, reader MessageReader) *Scribe {
	return &Scribe{db, ctx, reader}
}

func asTime(s discordgo.Timestamp) time.Time {
	t, _ := time.Parse(time.RFC3339, string(s))
	return t
}

func (s *Scribe) GetMessages(channelID string, guildID string) {
	messages, err := s.reader.ChannelMessages(channelID, 100, "", "", "")
	if err != nil {
		log.Fatal(err)
	}

	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}

	tx.Exec("SET NAMES utf8mb4;")
	for _, message := range messages {
		authorID := message.Author.ID
		fmt.Println("Reading message", message.ID, message.ChannelID, message.GuildID, authorID, message.Content, message.Timestamp)
		_, exec_err := tx.Exec(`INSERT INTO
		messages (id, channel_id, guild_id, author_id, content, timestamp)
		values (?, ?, ?, ?, ?, ?)`,
			message.ID,
			message.ChannelID,
			guildID,
			authorID,
			message.Content,
			asTime(message.Timestamp))

		if exec_err != nil {
			log.Fatalf("Error adding message: %v", exec_err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

}

func (s *Scribe) AddMessage(message *discordgo.Message) {
	s.GetMessages(message.ChannelID, message.GuildID)
}
