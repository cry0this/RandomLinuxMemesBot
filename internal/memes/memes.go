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

func GetRandomMemeURL(ctx context.Context, guildID string) (string, error) {
	client := gohttpclient.New(memeURL)

	response, err := client.Get(ctx, "/gimme/linuxmemes/50")
	if err != nil {
		return "", err
	}

	var data MemeResponse
	if err := response.Unmarshal(&data); err != nil {
		return "", err
	}

	var url string

	for _, m := range data.Memes {
		log := logrus.WithFields(logrus.Fields{
			"url":     m.URL,
			"guildID": guildID,
		})

		if m.NSFW && !nsfwEnabled {
			log.Info("meme url skipped because NSFW")
			continue
		}

		exist, err := redis.ExistsInCache(ctx, guildID, m.URL)
		if err != nil {
			return "", err
		}

		if exist {
			log.Info("url was in cache, skipping")
			continue
		}

		url = m.URL
		break
	}

	if url == "" {
		return "", fmt.Errorf("unable to find new meme :(")
	}

	err = redis.AddToCache(ctx, guildID, url)
	if err != nil {
		return "", err
	}

	return url, nil
}
