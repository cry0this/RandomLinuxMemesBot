package memes

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/sethvargo/go-retry"
	"github.com/sirupsen/logrus"

	"github.com/cry0this/RandomLinuxMemesBot/internal/reddit"
	"github.com/cry0this/RandomLinuxMemesBot/internal/redis"
)

const (
	maxTries = 5
)

func Init(ctx context.Context) error {
	return fillCache(ctx)
}

func GetNewMeme(ctx context.Context, guildId string) (*reddit.Post, error) {
	log := logrus.WithField("guild", guildId)
	log.Info("got new meme request")

	cachedAll, err := redis.GetCachedPosts(ctx, "")
	if err != nil {
		return nil, err
	}

	if len(cachedAll) == 0 {
		if err := fillCache(ctx); err != nil {
			return nil, err
		}
		cachedAll, err = redis.GetCachedPosts(ctx, "")
		if err != nil {
			return nil, err
		}
	}

	cachedGuild, err := redis.GetCachedPosts(ctx, guildId)
	if err != nil {
		return nil, err
	}

	// randomize posts
	randomized := make([]*reddit.Post, len(cachedAll))
	copy(randomized, cachedAll)
	rand.Shuffle(len(randomized), func(i, j int) { randomized[i], randomized[j] = randomized[j], randomized[i] })

	for _, p := range randomized {
		if !postInSlice(p, cachedGuild) {
			log.WithFields(logrus.Fields{
				"ID":    p.ID,
				"Title": p.Title,
				"URL":   p.URL,
			}).Info("adding post to guild cache")

			if err := redis.AddToCache(ctx, p, guildId); err != nil {
				return nil, err
			}

			return p, nil
		}
	}

	// try to find newer posts
	before := cachedAll[0].ID
	postsBefore, err := reddit.GetPostsBefore(ctx, before)
	if err != nil {
		return nil, err
	}

	if len(postsBefore) > 0 {
		log.Info("adding newer posts to global cache")

		if err := redis.PushToHead(ctx, postsBefore, ""); err != nil {
			return nil, err
		}

		p := postsBefore[0]

		log.WithFields(logrus.Fields{
			"ID":    p.ID,
			"Title": p.Title,
			"URL":   p.URL,
		}).Info("adding post to guild cache")

		if err := redis.AddToCache(ctx, p, guildId); err != nil {
			return nil, err
		}
		return p, nil
	}

	postsAfter := make([]*reddit.Post, 0)
	var newPost *reddit.Post

	// try to find older posts
	b := retry.NewConstant(time.Millisecond)
	b = retry.WithMaxRetries(maxTries, b)
	if retry.Do(ctx, b, func(_ context.Context) error {
		log.Info("trying to find older post for guild...")

		after := cachedAll[len(cachedAll)-1].ID
		postsAfter, err = reddit.GetPostsAfter(ctx, after)
		if err != nil {
			return err
		}

		if len(postsAfter) == 0 {
			return errors.New("unable to find older posts")
		}

		log.Info("adding older posts to cache")
		if err := redis.PushToTail(ctx, postsAfter, ""); err != nil {
			return err
		}

		cachedAll, err = redis.GetCachedPosts(ctx, "")
		if err != nil {
			return err
		}

		for _, p := range postsAfter {
			if !postInSlice(p, cachedGuild) {
				newPost = p
				return nil
			}
		}

		return retry.RetryableError(errors.New("unable to find older post for guild"))
	}); err != nil {
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"ID":    newPost.ID,
		"Title": newPost.Title,
		"URL":   newPost.URL,
	}).Info("adding post to guild cache")

	if err := redis.AddToCache(ctx, newPost, guildId); err != nil {
		return nil, err
	}

	return newPost, nil
}

func fillCache(ctx context.Context) error {
	p, err := redis.GetCachedPosts(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to read cache: %v", err)
	}

	if len(p) == 0 {
		p, err = reddit.GetPosts(ctx)
		if err != nil {
			return fmt.Errorf("failed to get reddit posts: %v", err)
		}

		if err := redis.PushToHead(ctx, p, ""); err != nil {
			return fmt.Errorf("failed to fill cache: %v", err)
		}
	}

	return nil
}

func postInSlice(p *reddit.Post, posts []*reddit.Post) bool {
	for _, i := range posts {
		if *p == *i {
			return true
		}
	}
	return false
}
