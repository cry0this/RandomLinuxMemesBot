package discord

import (
	"github.com/bwmarrin/discordgo"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
	"github.com/cry0this/RandomLinuxMemesBot/internal/memes"
)

func getTopMeme(s *discordgo.Session, i *discordgo.InteractionCreate) {
	actx := appctx.NewContext(*ctx)

	actx.Logger.WithFields(getLogFields(i)).Info("invoked new command")
	log := actx.Logger

	ID := getID(i)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})

	options := i.ApplicationCommandData().Options
	optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optMap[opt.Name] = opt
	}

	period, ok := optMap["period"]
	if !ok {
		log.Error("failed to get 'period' option")
		followUpErrMessage(actx, s, i, ErrorMsg{
			Default: "failed to get period",
			Localizations: map[discordgo.Locale]string{
				discordgo.Russian: "не удалось получить период",
			},
		})
	}

	errMsg := ErrorMsg{
		Default: "Ooops! Couldn't find new meme :(\nTry another period or run later...",
		Localizations: map[discordgo.Locale]string{
			discordgo.Russian: "Упс! Не получилось найти новый мем :(\nПопробуй другой период или запустить позже...",
		},
	}

	meme, err := memes.GetTopMeme(actx, ID, period.StringValue())
	if err != nil {
		log.WithError(err).Error("failed to get meme url")
		followUpErrMessage(actx, s, i, errMsg)
		return
	}

	postMeme(actx, s, i, meme, errMsg)
	log.Info("done")
}
