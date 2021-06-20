package scribe

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Intentionally shaped to fit discordgo.Session, but left flexible in case
// someone wants to fake/stub it.
type MessageReader interface {
	ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) (st []*discordgo.Message, err error)
}

type Scribe struct {
	db     *sql.DB
	ctx    context.Context
	reader MessageReader
}

func NewScribe(db *sql.DB, ctx context.Context, reader MessageReader) *Scribe {
	return &Scribe{db, ctx, reader}
}

type SlurpMessagesResult struct {
	Complete         bool // true if all messages are in the database
	MessagesRead     int  // number of messages read from Discord
	BeginningReached bool // true if the beginning of the channel has been reached during this read
	ReadingEarlier   bool // true if was reading earlier messages
}

// Calling function will read at most 100 messages from Discord and put them
// into the `messages` table in the SQL database. This function is a bit clever
// in how it reads these messages. The general approach is to first read earliest
// messages until the very beginning is reached. Then, read the more recent
// messages until there's no more to read.
func (s *Scribe) SlurpMessages(channelID string, guildID string) (*SlurpMessagesResult, error) {
	result := &SlurpMessagesResult{false, 0, false, false}
	// Wrap all this work in one transaction.
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return result, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	q := &query{tx}

	r, err := q.getChannelRange(channelID)
	if err != nil {
		return result, fmt.Errorf("failed to get range for channel '%v': %v", channelID, err)
	}

	hasBeginning, err := q.getChannelHasBeginning(channelID)
	if err != nil {
		return result, fmt.Errorf("faild to get has_beginning for channel '%v': %v", channelID, err)
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
		result.ReadingEarlier = true
	} else {
		afterID = r.latestID
	}

	messages, err := s.reader.ChannelMessages(channelID, 100, beforeID, afterID, "")
	if err != nil {
		return result, fmt.Errorf("failed to read messages from Discord for channel '%v': %v", channelID, err)
	}

	result.MessagesRead = len(messages)

	q.storeMessages(channelID, guildID, messages)

	// If the result contains fewer than 100 earler messages,
	// presume that the beginning has been reached.
	if len(messages) < 100 {
		if !hasBeginning {
			result.BeginningReached = true
			if err := q.setChannelHasBeginning(channelID); err != nil {
				return result, fmt.Errorf("failed to update has_beginning for channel '%v': %v", channelID, err)
			}
		} else {
			result.Complete = true
		}
	}

	// Let's go!
	if err := tx.Commit(); err != nil {
		return result, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return result, nil
}
