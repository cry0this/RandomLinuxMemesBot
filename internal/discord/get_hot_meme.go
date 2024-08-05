package discord

import (
	"github.com/bwmarrin/discordgo"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
	"github.com/cry0this/RandomLinuxMemesBot/internal/memes"
)

func getHotMeme(s *discordgo.Session, i *discordgo.InteractionCreate) {
	actx := appctx.NewContext(*ctx)

	actx.Logger.WithFields(getLogFields(i)).Info("invoked new command")
	log := actx.Logger

	ID := getID(i)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})

	errMsg := ErrorMsg{
		Default: "Ooops! Couldn't find new meme :(\nTry again later...",
		Localizations: map[discordgo.Locale]string{
			discordgo.Russian: "Упс! Не получилось найти новый мем :(\nПопробуй позже...",
		},
	}

	meme, err := memes.GetHotMeme(actx, ID)
	if err != nil {
		log.WithError(err).Error("failed to get meme url")
		followUpErrMessage(actx, s, i, errMsg)
		return
	}

	postMeme(actx, s, i, meme, errMsg)
	log.Info("done")
}
