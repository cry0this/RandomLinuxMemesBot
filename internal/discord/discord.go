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

func Init(c context.Context) error {
	ctx = &c

	log := logrus.WithField("module", "discord")
	log.Info("initializing...")

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return errors.New("empty discord token, check DISCORD_TOKEN variable")
	}

	log.Info("checking token...")

	var err error
	session, err = discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("unable to initialize discord bot: %v", err)
	}

	log.Info("setting up...")

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := cmdHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.WithFields(logrus.Fields{
			"username":      s.State.User.Username,
			"discriminator": s.State.User.Discriminator,
		}).Info("logged in")
	})

	log.Info("opening session...")

	err = session.Open()
	if err != nil {
		return fmt.Errorf("cannot open discord session: %v", err)
	}

	log.Info("adding commands...")
	for _, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", v)
		if err != nil {
			return fmt.Errorf("cannot add discord '%s' command: %v", v.Name, err)
		}
		registeredCommands = append(registeredCommands, cmd)
	}

	log.Info("initialized")
	return nil
}

func Cleanup() error {
	log := logrus.WithField("module", "discord")
	log.Info("cleaning up...")

	log.Info("deleting commands...")
	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, "", v.ID)
		if err != nil {
			return fmt.Errorf("cannot delete discord '%s' command: %v", v.Name, err)
		}
	}

	log.Info("shutting down...")
	if err := session.Close(); err != nil {
		return err
	}

	log.Info("cleaned up")

	return nil
}
