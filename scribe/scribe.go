package scribe

import (
	"context"
	"database/sql"
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

	// Wrap all this work in one transaction.
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// Check to see if we know about this channel.
	var (
		isFullyRead bool
		beforeID    string
		sqlBeforeID sql.NullString
	)
	row := tx.QueryRow(`
		SELECT is_fully_read, earliest_read_message 
		FROM channels WHERE id = ?`, channelID)
	if err := row.Scan(&isFullyRead, &sqlBeforeID); err != nil {
		if err == sql.ErrNoRows {
			// we don't know about this channel, let's learn about it!
			isFullyRead = false
			beforeID = ""
			if _, exec_err := tx.Exec(`
				INSERT INTO channels (id, is_fully_read)
				VALUES (?, ?)`, channelID, false); exec_err != nil {
				log.Fatal(exec_err)
			}
		} else {
			log.Fatal(err)
		}
	} else {
		beforeID = sqlBeforeID.String
	}

	if !isFullyRead {

		messages, err := s.reader.ChannelMessages(channelID, 100, beforeID, "", "")
		if err != nil {
			log.Fatal(err)
		}

		tx.Exec("SET NAMES utf8mb4;") // make emoji be storable.

		stmt, err := tx.Prepare(`
		INSERT INTO messages (id, channel_id, guild_id, author_id, content, timestamp)
		VALUES (?, ?, ?, ?, ?, ?)`)

		if err != nil {
			log.Fatal(err)
		}

		for _, message := range messages {
			authorID := message.Author.ID
			_, exec_err := stmt.Exec(
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

		beforeID = messages[len(messages)-1].ID
		// TODO: consider the case when exactly 100 messages are read and 100th is the last one.
		isFullyRead = len(messages) < 100

		if _, err := tx.Exec(`
		UPDATE channels SET is_fully_read = ?, earliest_read_message = ?`, isFullyRead, beforeID); err != nil {
			log.Fatal(err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}

	}

}

func (s *Scribe) AddMessage(message *discordgo.Message) {
	s.GetMessages(message.ChannelID, message.GuildID)
}
