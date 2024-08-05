package memes

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/sethvargo/go-retry"
	"github.com/sirupsen/logrus"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
	"github.com/cry0this/RandomLinuxMemesBot/internal/reddit"
	"github.com/cry0this/RandomLinuxMemesBot/internal/redis"
)

const (
	maxTries = 50
)

func GetHotMeme(actx *appctx.Context, guild string) (*reddit.Post, error) {
	log := actx.Logger.WithFields(logrus.Fields{
		"func":  "memes.GetHotMeme",
		"guild": guild,
	})
	log.Info("got new request")

	cachedAll, err := redis.GetCachedPosts(actx, "")
	if err != nil {
		return nil, err
	}

	if len(cachedAll) == 0 {
		log.Info("cache is empty")
		if err := fillCache(actx); err != nil {
			return nil, err
		}
		cachedAll, err = redis.GetCachedPosts(actx, "")
		if err != nil {
			return nil, err
		}
	}

	cachedGuild, err := redis.GetCachedPosts(actx, guild)
	if err != nil {
		return nil, err
	}

	// randomize posts
	randomized := make([]*reddit.Post, len(cachedAll))
	copy(randomized, cachedAll)
	rand.Shuffle(len(randomized), func(i, j int) { randomized[i], randomized[j] = randomized[j], randomized[i] })

	for _, p := range randomized {
		if !postInSlice(p, cachedGuild) {
			log.WithField("post", p).Info("adding post to guild cache")
			if err := redis.AddToCache(actx, p, guild); err != nil {
				return nil, err
			}
			return p, nil
		}
	}

	log.Info("trying to find new posts...")

	page := reddit.NewHotPage(actx)
	page.Posts = append(page.Posts, cachedAll[0])

	page, err = page.PrevPage()
	if err == nil {
		postsBefore, err := page.GetPosts()
		if err == nil {
			postsBefore = reddit.FilterPosts(actx, postsBefore)
			if len(postsBefore) > 0 {
				log.Info("adding newer posts to global cache")
				if err := redis.PushToHead(actx, postsBefore, ""); err != nil {
					return nil, err
				}

				p := postsBefore[0]

				log.WithField("post", p).Info("adding post to guild cache")
				if err := redis.AddToCache(actx, p, guild); err != nil {
					return nil, err
				}

				return p, nil
			}
		}
	}

	page = reddit.NewHotPage(actx)
	page.Posts = append(page.Posts, cachedAll[0])

	var newPost *reddit.Post

	// try to find older posts
	b := retry.NewConstant(time.Millisecond)
	b = retry.WithMaxRetries(maxTries, b)
	if err := retry.Do(actx.Context, b, func(_ context.Context) error {
		log.Info("trying to find older post for guild...")

		page, err = page.NextPage()
		if err != nil {
			return err
		}

		postsAfter, err := page.GetPosts()
		if err != nil {
			return err
		}

		postsAfter = reddit.FilterPosts(actx, postsAfter)

		if len(postsAfter) == 0 {
			return retry.RetryableError(errors.New("unable to find older posts"))
		}

		log.Info("adding older posts to cache")
		if err := redis.PushToTail(actx, postsAfter, ""); err != nil {
			return err
		}

		cachedAll, err = redis.GetCachedPosts(actx, "")
		if err != nil {
			return err
		}

		for _, p := range postsAfter {
			if !postInSlice(p, cachedGuild) {
				newPost = p
				return nil
			}
		}

		return retry.RetryableError(errors.New("unable to find older post"))
	}); err != nil {
		return nil, err
	}

	log.WithField("post", newPost).Info("adding post to guild cache")

	if err := redis.AddToCache(actx, newPost, guild); err != nil {
		return nil, err
	}

	return newPost, nil
}

func GetTopMeme(actx *appctx.Context, guild string, period string) (*reddit.Post, error) {
	log := actx.Logger.WithFields(logrus.Fields{
		"func":   "memes.GetTopMeme",
		"guild":  guild,
		"period": period,
	})

	log.Info("got new request")

	cache, err := redis.GetCachedPosts(actx, guild)
	if err != nil {
		return nil, err
	}

	var post *reddit.Post
	page := reddit.NewTopPage(actx, period)

	b := retry.NewConstant(time.Millisecond)
	b = retry.WithMaxRetries(maxTries, b)
	if err := retry.Do(actx.Context, b, func(_ context.Context) error {
		log.Info("trying to find new post")

		posts, err := page.GetPosts()
		if err != nil {
			return err
		}

		posts = reddit.FilterPosts(actx, posts)
		if len(posts) == 0 {
			page, err = page.NextPage()
			if err != nil {
				return err
			}
			return retry.RetryableError(errors.New("unable to find new post"))
		}

		rand.Shuffle(len(posts), func(i, j int) { posts[i], posts[j] = posts[j], posts[i] })

		for _, p := range posts {
			if !postInSlice(p, cache) {
				post = p
				return nil
			}
		}

		page, err = page.NextPage()
		if err != nil {
			return err
		}

		return retry.RetryableError(errors.New("unable to find new post"))
	}); err != nil {
		return nil, err
	}

	log.WithField("post", post).Info("adding post to guild cache")
	redis.AddToCache(actx, post, guild)

	return post, nil
}
