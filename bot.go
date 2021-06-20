package main

import (
	"fmt"
	"log"
	"time"

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

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

// discordgo callback: called after the when new message is posted.
func (b *bot) messageCreate(s *discordgo.Session, event *discordgo.MessageCreate) {
	defer duration(track("Add Message"))
	message := event.Message
	r, err := b.scribe.SlurpMessages(message.ChannelID, message.GuildID)
	if err != nil {
		log.Fatalf("failed to slurp messages: %v", err)
	}
	log.Println("Succesfully slurped messages:")
	log.Printf("\tcomplete: %v\n", r.Complete)
	log.Printf("\ttmessages read: %v\n", r.MessagesRead)
	log.Printf("\tbeginning reached: %v\n", r.BeginningReached)
	log.Printf("\tReading earlier messages: %v\n", r.ReadingEarlier)
}
