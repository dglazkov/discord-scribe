package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type bot struct {
	session    *discordgo.Session
	controller Controller
}

func newBot(s *discordgo.Session, c Controller) *bot {
	result := &bot{
		session:    s,
		controller: c,
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
	fmt.Println("Message created")
}
