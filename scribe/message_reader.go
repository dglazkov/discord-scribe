package scribe

import "github.com/bwmarrin/discordgo"

type MessageReader interface {
	ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) (st []*discordgo.Message, err error)
}
