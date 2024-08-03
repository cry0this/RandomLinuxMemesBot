package memes

import (
	"context"
	"fmt"
	"os"

	gohttpclient "github.com/bozd4g/go-http-client"
	"github.com/sirupsen/logrus"

	"github.com/cry0this/RandomLinuxMemesBot/internal/redis"
)

var memeURL string
var nsfwEnabled bool

type Meme struct {
	PostLink  string   `json:"postLink"`
	Subreddit string   `json:"subreddit"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	NSFW      bool     `json:"nsfw"`
	Spoiler   bool     `json:"spoiler"`
	Author    string   `json:"author"`
	Ups       int      `json:"ups"`
	Preview   []string `json:"preview"`
}

type MemeResponse struct {
	Count int    `json:"count"`
	Memes []Meme `json:"memes"`
}

func Init(nsfw bool) {
	logrus.WithField("NSFW", nsfw).Info("initializing meme api...")

	memeURL = os.Getenv("MEME_API_URL")
	if memeURL == "" {
		logrus.Fatal("failed to parse meme api url")
	}

	nsfwEnabled = nsfw
}

func GetRandomMeme(ctx context.Context, guildID string) (*Meme, error) {
	client := gohttpclient.New(memeURL)

	response, err := client.Get(ctx, "/gimme/linuxmemes/50")
	if err != nil {
		return nil, err
	}

	var data MemeResponse
	if err := response.Unmarshal(&data); err != nil {
		return nil, err
	}

	var meme *Meme

	for _, m := range data.Memes {
		log := logrus.WithFields(logrus.Fields{
			"url":     m.URL,
			"guildID": guildID,
		})

		if m.NSFW && !nsfwEnabled {
			log.Info("meme url skipped because NSFW")
			continue
		}

		exist, err := redis.IsCached(ctx, guildID, m.URL)
		if err != nil {
			return nil, err
		}

		if exist {
			log.Info("url was in cache, skipping")
			continue
		}

		meme = &m
		break
	}

	if meme == nil {
		return nil, fmt.Errorf("unable to find new meme :(")
	}

	err = redis.AddToCache(ctx, guildID, meme.URL)
	if err != nil {
		return nil, err
	}

	return meme, nil
}
