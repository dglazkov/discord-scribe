package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dglazkov/discord-scribe/scribe"
)

type bot struct {
	session *discordgo.Session
	scribe  *scribe.Scribe
}

func newBot(s *discordgo.Session, scribe *scribe.Scribe) *bot {
	result := &bot{
		session: s,
		scribe:  scribe,
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

// discordgo callback: called after the when new message is posted.
func (b *bot) messageCreate(s *discordgo.Session, event *discordgo.MessageCreate) {
	b.scribe.AddMessage(event.Message)
}
