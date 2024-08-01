package discord

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/gabriel-vasile/mimetype"
	"github.com/sirupsen/logrus"

	"github.com/cry0this/RandomLinuxMemesBot/internal/memes"
	"github.com/cry0this/RandomLinuxMemesBot/internal/utils"
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
	"get-random-linux-meme": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var fields logrus.Fields
		var ID string

		if i.User != nil {
			fields = logrus.Fields{
				"Command": "get-random-linux-meme",
				"User":    i.User.Username,
				"ID":      i.User.ID,
			}

			ID = i.User.ID
		}

		if i.Member != nil {
			fields = logrus.Fields{
				"Command":   "get-random-linux-meme",
				"GuildID":   i.GuildID,
				"ChannelID": i.ChannelID,
				"Member":    i.Member.User.Username,
			}

			ID = i.GuildID
		}

		log := logrus.WithFields(fields)
		log.Info("invoked new command")

		url, err := memes.GetRandomMemeURL(*ctx, ID)
		if err != nil {
			log.WithError(err).Error("failed to get meme url")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ooops! Couldn't load new linux meme :(",
				},
			})
		}

		log.Infof("got meme url: %s", url)

		file, err := os.CreateTemp("/tmp", "linuxmemes.*.jpg")
		if err != nil {
			log.WithError(err).Error("failed to create tmp file")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ooops! Couldn't load new linux meme :(",
				},
			})
		}

		defer os.Remove(file.Name())

		log = log.WithFields(logrus.Fields{
			"url":  url,
			"file": file.Name(),
		})

		log.Info("downloading meme file...")

		err = utils.DownloadFile(file.Name(), url)
		if err != nil {
			log.WithError(err).Error("failed to download meme file")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ooops! Couldn't load new linux meme :(",
				},
			})
		}

		log.Info("detecting mime type...")

		mime, err := mimetype.DetectFile(file.Name())
		if err != nil {
			log.WithError(err).Error("failed to detect mime type")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ooops! Couldn't load new linux meme :(",
				},
			})
		}

		log.Infof("mime type: %s", mime.String())

		reader, err := os.Open(file.Name())
		if err != nil {
			log.WithError(err).Error("failed to open file")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ooops! Couldn't load new linux meme :(",
				},
			})
		}

		log.Info("uploading file to discord")

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Enjoy fresh linux meme :)",
				Files: []*discordgo.File{
					{
						ContentType: mime.String(),
						Name:        file.Name(),
						Reader:      reader,
					},
				},
			},
		})

		log.Info("done")
	},
}

var registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
