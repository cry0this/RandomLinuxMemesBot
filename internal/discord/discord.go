package discord

import (
	"context"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	ctx                *context.Context
	session            *discordgo.Session
	registeredCommands []*discordgo.ApplicationCommand
)

func Init(c *context.Context) {
	ctx = c

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		logrus.Fatal("unable to get discord token from env")
	}

	logrus.Info("checking discord token...")

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		logrus.WithError(err).Fatal("unable to initialize discord bot")
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
		logrus.WithError(err).Fatal("cannot open discord session")
	}

	logrus.Info("adding discord commands...")
	for _, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", v)
		if err != nil {
			logrus.WithField("command", v.Name).WithError(err).Fatal("cannot create discord command")
		}
		registeredCommands = append(registeredCommands, cmd)
	}

	logrus.Info("discord initialized")
}

func Cleanup() {
	logrus.Info("removing discord commands...")

	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, "", v.ID)
		if err != nil {
			logrus.WithField("command", v.Name).WithError(err).Fatal("cannot delete discord command")
		}
	}

	logrus.Info("gracefully shutting down")
	session.Close()
}

func followUpErrMessage(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	var fields logrus.Fields

	if i.User != nil {
		fields = logrus.Fields{
			"Command": i.ApplicationCommandData().Name,
			"User":    i.User.Username,
			"ID":      i.User.ID,
		}
	}

	if i.Member != nil {
		fields = logrus.Fields{
			"Command":   i.ApplicationCommandData().Name,
			"GuildID":   i.GuildID,
			"ChannelID": i.ChannelID,
			"Member":    i.Member.User.Username,
		}
	}

	s.InteractionResponseDelete(i.Interaction)

	if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: msg,
		Flags:   discordgo.MessageFlagsEphemeral,
	}); err != nil {
		logrus.WithFields(fields).WithError(err).Error("failed to follow up message")
	}

}
