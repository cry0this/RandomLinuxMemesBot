package discord

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
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
