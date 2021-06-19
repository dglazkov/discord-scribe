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

func (s *Scribe) GetMessages(channelID string, guildID string) {
	// Wrap all this work in one transaction.
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	q := &query{tx}

	r, err := q.getChannelRange(channelID)
	if err != nil {
		log.Fatal(err)
	}

	hasBeginning, err := q.getChannelHasBeginning(channelID)
	if err != nil {
		log.Fatal(err)
	}

	var (
		beforeID string = ""
		afterID  string = ""
	)
	// Strategy:
	// 	if the store doesn't contain messages from the very beginning,
	// 	keep reading earlier messages (messages that came before the
	// 	earliest known message).
	// 	Otherwise, read later messages (messages that came after
	// 	latest known message).
	if !hasBeginning {
		beforeID = r.earliestID
	} else {
		afterID = r.latestID
	}

	messages, err := s.reader.ChannelMessages(channelID, 100, beforeID, afterID, "")
	if err != nil {
		log.Fatal(err)
	}

	q.storeMessages(channelID, guildID, messages)

	// If the result contains fewer than 100 earler messages,
	// presume that the beginning has been reached.
	if !hasBeginning && len(messages) < 100 {
		if err := q.setChannelHasBeginning(channelID); err != nil {
			log.Fatal(err)
		}
	}

	// Let's go!
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

func (s *Scribe) AddMessage(message *discordgo.Message) {
	defer duration(track("Add Message"))
	s.GetMessages(message.ChannelID, message.GuildID)
}
