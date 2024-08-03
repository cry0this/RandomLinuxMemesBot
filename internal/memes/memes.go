package memes

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/cry0this/RandomLinuxMemesBot/internal/reddit"
	"github.com/cry0this/RandomLinuxMemesBot/internal/redis"
	"github.com/sirupsen/logrus"
)

const ttl = 6 * time.Hour

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

		if err := redis.AddToCache(ctx, p, guildId); err != nil {
			return nil, err
		}
		return p, nil
	}

	// try to find older posts
	after := cachedAll[len(cachedAll)-1].ID
	postsAfter, err := reddit.GetPostsAfter(ctx, after)
	if err != nil {
		return nil, err
	}

	if len(postsAfter) == 0 {
		return nil, errors.New("unable to find new posts")
	}

	log.Info("adding older posts to cache")

	if err := redis.PushToTail(ctx, postsAfter, ""); err != nil {
		return nil, err
	}

	p := postsAfter[0]
	if err := redis.AddToCache(ctx, p, guildId); err != nil {
		return nil, err
	}

	return p, nil
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

	if err := redis.Expire(ctx, "", ttl); err != nil {
		return err
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
