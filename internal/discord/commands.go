package discord

import (
	"github.com/bwmarrin/discordgo"
)

var permissions int64 = discordgo.PermissionSendMessages & discordgo.PermissionAttachFiles

var commands = []*discordgo.ApplicationCommand{
	{
		Name:                     "get-random-linux-meme",
		Description:              "Get random linux meme",
		DefaultMemberPermissions: &permissions,
	},
}

var cmdHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"get-random-linux-meme": getRandomLinuxMeme,
}
