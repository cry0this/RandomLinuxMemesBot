package discord

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/gabriel-vasile/mimetype"
	"github.com/sirupsen/logrus"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
	"github.com/cry0this/RandomLinuxMemesBot/internal/reddit"
)

func downloadFile(fileName string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	if written == 0 {
		return errors.New("empty file")
	}

	return nil
}

func getLogFields(i *discordgo.InteractionCreate) logrus.Fields {
	var fields logrus.Fields

	if i.User != nil {
		fields = logrus.Fields{
			"command": i.ApplicationCommandData().Name,
			"user":    i.User.Username,
			"id":      i.User.ID,
		}
	}

	if i.Member != nil {
		fields = logrus.Fields{
			"command": i.ApplicationCommandData().Name,
			"guild":   i.GuildID,
			"channel": i.ChannelID,
			"member":  i.Member.User.Username,
		}
	}

	return fields
}

func getID(i *discordgo.InteractionCreate) string {
	var ID string

	if i.User != nil {
		ID = i.User.ID
	}

	if i.Member != nil {
		ID = i.GuildID
	}

	return ID
}

func followUpErrMessage(actx *appctx.Context, s *discordgo.Session, i *discordgo.InteractionCreate, msg ErrorMsg) {
	s.InteractionResponseDelete(i.Interaction)

	if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: msg.Get(i.Locale),
		Flags:   discordgo.MessageFlagsEphemeral,
	}); err != nil {
		actx.Logger.WithField("func", "discord.followUpErrMessage").WithError(err).Error("failed to follow up message")
	}

}

func postMeme(actx *appctx.Context, s *discordgo.Session, i *discordgo.InteractionCreate, meme *reddit.Post, errMsg ErrorMsg) {
	log := actx.Logger
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
		log.WithError(err).Error("failed to respond")
	}
}
