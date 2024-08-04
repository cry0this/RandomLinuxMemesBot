package discord

import (
	"context"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
	"github.com/cry0this/RandomLinuxMemesBot/internal/memes"
)

func getRandomLinuxMeme(s *discordgo.Session, i *discordgo.InteractionCreate) {
	u := uuid.New()
	c := context.WithValue(*ctx, appctx.Identifier, u.String())
	actx := appctx.NewContext(c)

	fields := getLogFields(i)
	actx.Logger.WithFields(fields).Info("invoked new command")
	log := actx.Logger

	ID := getID(i)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})

	errMsg := "Ooops! Couldn't find new linux meme :(\nTry again later..."

	meme, err := memes.GetNewMeme(actx, ID)
	if err != nil {
		log.WithError(err).Error("failed to get meme url")
		followUpErrMessage(actx, s, i, errMsg)
		return
	}

	log.Infof("got meme: %v", meme)

	file, err := os.CreateTemp("/tmp", "linuxmemes.*.jpg")
	if err != nil {
		log.WithError(err).Error("failed to create tmp file")
		followUpErrMessage(actx, s, i, errMsg)
		return
	}

	defer os.Remove(file.Name())

	log.WithFields(logrus.Fields{
		"url":  meme.URL,
		"file": file.Name(),
	}).Info("downloading meme file...")

	err = downloadFile(file.Name(), meme.URL)
	if err != nil {
		log.WithError(err).Error("failed to download meme file")
		followUpErrMessage(actx, s, i, errMsg)
		return
	}

	log.Info("detecting mime type...")

	mime, err := mimetype.DetectFile(file.Name())
	if err != nil {
		log.WithError(err).Error("failed to detect mime type")
		followUpErrMessage(actx, s, i, errMsg)
		return
	}

	log.Infof("mime type: %s", mime.String())

	reader, err := os.Open(file.Name())
	if err != nil {
		log.WithError(err).Error("failed to open file")
		followUpErrMessage(actx, s, i, errMsg)
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
