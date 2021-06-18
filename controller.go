package main

import "github.com/bwmarrin/discordgo"

type DiscordController struct {
	session *discordgo.Session
}

func (dc *DiscordController) ChannelMessage(channelID, messageID string) (st *discordgo.Message, err error) {
	return dc.session.ChannelMessage(channelID, messageID)
}

func (dc *DiscordController) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) (st []*discordgo.Message, err error) {
	return dc.session.ChannelMessages(channelID, limit, beforeID, afterID, aroundID)
}

func (dc *DiscordController) GuildChannelCreateComplex(guildID string, data discordgo.GuildChannelCreateData) (st *discordgo.Channel, err error) {
	return dc.session.GuildChannelCreateComplex(guildID, data)
}

func (dc *DiscordController) ChannelEditComplex(channelID string, data *discordgo.ChannelEdit) (st *discordgo.Channel, err error) {
	return dc.session.ChannelEditComplex(channelID, data)
}

func (dc *DiscordController) GuildChannelsReorder(guildID string, channels []*discordgo.Channel) error {
	return dc.session.GuildChannelsReorder(guildID, channels)
}
