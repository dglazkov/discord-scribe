package scribe

import (
	"database/sql"
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
)

type channelRange struct {
	earliestID string
	latestID   string
}

type query struct {
	tx *sql.Tx
}

func (q *query) getChannelRange(channelID string) (*channelRange, error) {
	rows, err := q.tx.Query(`
		SELECT id 
		FROM messages AS m, 
		(SELECT MIN(timestamp) AS earliest, MAX(timestamp) AS latest 
			FROM messages
			WHERE channel_id = ?) AS t
		WHERE 
			(t.earliest = m.timestamp OR t.latest = m.timestamp)`, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// get earliest id
	var sqlEarliestID sql.NullString
	if !rows.Next() {
		return &channelRange{}, nil
	}
	if err := rows.Scan(&sqlEarliestID); err != nil {
		return nil, err
	}

	// get latest id
	var sqlLatestID sql.NullString
	if !rows.Next() {
		// This is weird and shouldn't happen
		return nil, errors.New("somehow, the channel range returned only one row")
	}
	if err := rows.Scan(&sqlLatestID); err != nil {
		return nil, err
	}

	return &channelRange{sqlEarliestID.String, sqlLatestID.String}, nil
}

func (q *query) getChannelHasBeginning(channelID string) (bool, error) {
	var hasBeginning bool
	row := q.tx.QueryRow(`
		SELECT has_beginning 
		FROM channels WHERE id = ?`, channelID)
	if err := row.Scan(&hasBeginning); err != nil {
		if err == sql.ErrNoRows {
			// We don't know about this channel, let's learn about it!
			if _, err := q.tx.Exec(`
				INSERT INTO channels (id)
				VALUES (?)`, channelID); err != nil {
				return false, err
			}
		} else {
			return false, err
		}
	}
	return hasBeginning, nil
}

func (q *query) setChannelHasBeginning(channelID string) error {
	if _, err := q.tx.Exec(`
		UPDATE channels SET has_beginning = true
		WHERE id = ?`, channelID); err != nil {
		return err
	}
	return nil
}

func asTime(s discordgo.Timestamp) time.Time {
	t, _ := time.Parse(time.RFC3339, string(s))
	return t
}

func (q *query) storeMessages(channelID string, guildID string, messages []*discordgo.Message) error {
	q.tx.Exec("SET NAMES utf8mb4;") // Make emoji be storable.

	stmt, err := q.tx.Prepare(`
	INSERT INTO messages (id, channel_id, guild_id, author_id, content, timestamp)
	VALUES (?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return err
	}

	for _, message := range messages {
		authorID := message.Author.ID
		if _, err := stmt.Exec(
			message.ID,
			message.ChannelID,
			guildID,
			authorID,
			message.Content,
			asTime(message.Timestamp)); err != nil {
			return err
		}
	}

	return nil
}
