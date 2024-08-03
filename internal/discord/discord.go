package discord

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	ctx                *context.Context
	session            *discordgo.Session
	registeredCommands []*discordgo.ApplicationCommand
)

func Init(c *context.Context) error {
	ctx = c

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return errors.New("empty discord token, check DISCORD_TOKEN variable")
	}

	logrus.Info("checking discord token...")

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("unable to initialize discord bot: %v", err)
	}

	logrus.Info("initializing discord bot...")

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := cmdHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logrus.WithFields(logrus.Fields{
			"username":      s.State.User.Username,
			"discriminator": s.State.User.Discriminator,
		}).Info("Logged in")
	})

	logrus.Info("opening discord session...")

	err = session.Open()
	if err != nil {
		return fmt.Errorf("cannot open discord session: %v", err)
	}

	logrus.Info("adding discord commands...")
	for _, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", v)
		if err != nil {
			return fmt.Errorf("cannot register discord '%s' command: %v", v.Name, err)
		}
		registeredCommands = append(registeredCommands, cmd)
	}

	logrus.Info("discord initialized")
	return nil
}

func Cleanup() error {
	logrus.Info("removing discord commands...")

	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, "", v.ID)
		if err != nil {
			return fmt.Errorf("cannot delete discord '%s' command: %v", v.Name, err)
		}
	}

	logrus.Info("gracefully shutting down")
	session.Close()

	return nil
}

func followUpErrMessage(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	fields := getLogFields(i)
	log := logrus.WithFields(fields)

	s.InteractionResponseDelete(i.Interaction)

	if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: msg,
		Flags:   discordgo.MessageFlagsEphemeral,
	}); err != nil {
		log.WithError(err).Error("failed to follow up message")
	}

}
