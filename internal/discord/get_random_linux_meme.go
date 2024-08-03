package discord

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/gabriel-vasile/mimetype"
	"github.com/sirupsen/logrus"

	"github.com/cry0this/RandomLinuxMemesBot/internal/memes"
	"github.com/cry0this/RandomLinuxMemesBot/internal/utils"
)

func getRandomLinuxMeme(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})

	errMsg := "Ooops! Couldn't find new linux meme :("

	meme, err := memes.GetRandomMeme(*ctx, ID)
	if err != nil {
		log.WithError(err).Error("failed to get meme url")
		followUpErrMessage(s, i, errMsg)
		return
	}

	log.Infof("got meme: %v", meme)

	url := meme.Preview[len(meme.Preview)-1]

	file, err := os.CreateTemp("/tmp", "linuxmemes.*.jpg")
	if err != nil {
		log.WithError(err).Error("failed to create tmp file")
		followUpErrMessage(s, i, errMsg)
		return
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
		followUpErrMessage(s, i, errMsg)
		return
	}

	log.Info("detecting mime type...")

	mime, err := mimetype.DetectFile(file.Name())
	if err != nil {
		log.WithError(err).Error("failed to detect mime type")
		followUpErrMessage(s, i, errMsg)
		return
	}

	log.Infof("mime type: %s", mime.String())

	reader, err := os.Open(file.Name())
	if err != nil {
		log.WithError(err).Error("failed to open file")
		followUpErrMessage(s, i, errMsg)
		return
	}

	log.Info("uploading file to discord")

	if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("> %s", meme.Title),
		Files: []*discordgo.File{
			{
				ContentType: mime.String(),
				Name:        file.Name(),
				Reader:      reader,
			},
		},
	}); err != nil {
		log.WithError(err).Error("failed to response")
	}

	log.Info("done")
}
