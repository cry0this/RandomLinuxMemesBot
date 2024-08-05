package memes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sethvargo/go-retry"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
	"github.com/cry0this/RandomLinuxMemesBot/internal/reddit"
	"github.com/cry0this/RandomLinuxMemesBot/internal/redis"
)

func postInSlice(p *reddit.Post, posts []*reddit.Post) bool {
	for _, i := range posts {
		if *p == *i {
			return true
		}
	}
	return false
}

func fillCache(actx *appctx.Context) error {
	actx.Logger.WithField("func", "memes.fillCache").Info("refilling cache...")

	p, err := redis.GetCachedPosts(actx, "")
	if err != nil {
		return fmt.Errorf("failed to read cache: %v", err)
	}

	if len(p) > 0 {
		return nil
	}

	var posts []*reddit.Post
	page := reddit.NewHotPage(actx)

	b := retry.NewConstant(time.Millisecond)
	b = retry.WithMaxRetries(maxTries, b)
	if err := retry.Do(actx.Context, b, func(_ context.Context) error {
		var err error

		posts, err = page.GetPosts()
		if err != nil {
			return err
		}

		posts = reddit.FilterPosts(actx, posts)

		if len(posts) == 0 {
			page, err = page.NextPage()
			if err != nil {
				return err
			}
			return retry.RetryableError(errors.New("unable to find new posts"))
		}

		return nil
	}); err != nil {
		return err
	}

	if err := redis.PushToHead(actx, posts, ""); err != nil {
		return fmt.Errorf("failed to fill cache: %v", err)
	}

	return nil
}
