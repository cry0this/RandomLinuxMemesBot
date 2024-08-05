package discord

import (
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "hot-meme",
		Description: "Post random hot meme",
		NameLocalizations: &map[discordgo.Locale]string{
			discordgo.Russian: "свежий-мем",
		},
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.Russian: "Выложить случайный свежий мем",
		},
	},
	{
		Name:        "top-meme",
		Description: "Post random meme from top of some period",
		NameLocalizations: &map[discordgo.Locale]string{
			discordgo.Russian: "топовый-мем",
		},
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.Russian: "Выложить случайный мем из топа за некоторый период",
		},
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "period",
				Description: "Period for search",
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.Russian: "период",
				},
				DescriptionLocalizations: map[discordgo.Locale]string{
					discordgo.Russian: "Период для поиска",
				},
				Type: discordgo.ApplicationCommandOptionString,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name: "hour",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.Russian: "час",
						},
						Value: "hour",
					},
					{
						Name: "day",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.Russian: "день",
						},
						Value: "day",
					},
					{
						Name: "week",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.Russian: "неделя",
						},
						Value: "week",
					},
					{
						Name: "month",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.Russian: "месяц",
						},
						Value: "month",
					},
					{
						Name: "year",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.Russian: "год",
						},
						Value: "year",
					},
					{
						Name: "all time",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.Russian: "всё время",
						},
						Value: "all",
					},
				},
				Required: true,
			},
		},
	},
}

var cmdHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"hot-meme": getHotMeme,
	"top-meme": getTopMeme,
}
