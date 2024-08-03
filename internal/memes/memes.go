package memes

import (
	"context"
	"fmt"
	"os"
	"time"

	gohttpclient "github.com/bozd4g/go-http-client"
	"github.com/cry0this/RandomLinuxMemesBot/internal/redis"
	retry "github.com/sethvargo/go-retry"
	"github.com/sirupsen/logrus"
)

const maxRetries = 5

var (
	memeURL     string
	nsfwEnabled bool
)

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

type MemeApiResponse struct {
	Count int
	Memes []Meme
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
	b := retry.NewConstant(time.Millisecond)
	b = retry.WithMaxRetries(maxRetries, b)

	var meme *Meme

	if err := retry.Do(ctx, b, func(_ context.Context) error {
		client := gohttpclient.New(memeURL)
		response, err := client.Get(ctx, "/gimme/linuxmemes/50")
		if err != nil {
			return err
		}

		var api MemeApiResponse
		if err = response.Unmarshal(&api); err != nil {
			return err
		}

		for _, m := range api.Memes {
			log := logrus.WithFields(logrus.Fields{
				"guildID": guildID,
				"url":     m.URL,
			})

			if m.NSFW && !nsfwEnabled {
				log.Warning("meme skipped because NSFW")
				continue
			}

			exist, err := redis.IsCached(ctx, guildID, m.URL)
			if err != nil {
				return err
			}

			if exist {
				log.Warning("meme already shown")
				continue
			}

			meme = &m
		}

		if meme == nil {
			if _, err := client.Get(ctx, "/clear"); err != nil {
				return err
			}

			return retry.RetryableError(fmt.Errorf("couldn't find new meme for guild"))
		}

		redis.AddToCache(ctx, guildID, meme.URL)
		return nil
	}); err != nil {
		return nil, err
	}

	return meme, nil
}
