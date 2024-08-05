package discord

import (
	"github.com/bwmarrin/discordgo"
)

type ErrorMsg struct {
	Localizations map[discordgo.Locale]string
	Default       string
}

func (e ErrorMsg) Get(l discordgo.Locale) string {
	if m, ok := e.Localizations[l]; ok {
		return m
	}

	return e.Default
}
