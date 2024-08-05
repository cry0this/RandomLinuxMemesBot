package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
	"github.com/cry0this/RandomLinuxMemesBot/internal/reddit"
)

var client *redis.Client

func Init(ctx context.Context) error {
	log := logrus.WithField("module", "redis")
	log.Info("initializing...")

	url := os.Getenv("REDIS_URL")
	opts, err := redis.ParseURL(url)
	if err != nil {
		return fmt.Errorf("failed to parse REDIS_URL variable: %v", err)
	}

	log.WithField("url", url).Info("connecting...")

	client = redis.NewClient(opts)

	if client.Ping(ctx).String() != "ping: PONG" {
		return errors.New("failed to connect redis, check REDIS_URL variable")
	}

	log.Info("initialized")
	return nil
}

func AddToCache(actx *appctx.Context, post *reddit.Post, key string) error {
	actx.Logger.WithFields(logrus.Fields{
		"func": "redis.AddToCache",
		"post": post,
	}).Info("adding to cache...")

	posts := []*reddit.Post{post}
	return PushToTail(actx, posts, key)
}

func PushToHead(actx *appctx.Context, posts []*reddit.Post, key string) error {
	p, err := preparePosts(actx, posts)
	if err != nil {
		return err
	}

	k := normalizeKey(key)
	if err := client.LPush(actx.Context, k, p).Err(); err != nil {
		return err
	}

	actx.Logger.WithFields(logrus.Fields{
		"func": "redis.PushToHead",
		"key":  k,
	}).Infof("pushed: %d", len(posts))

	return nil
}

func PushToTail(actx *appctx.Context, posts []*reddit.Post, key string) error {
	p, err := preparePosts(actx, posts)
	if err != nil {
		return err
	}

	k := normalizeKey(key)
	if err := client.RPush(actx.Context, k, p).Err(); err != nil {
		return err
	}

	actx.Logger.WithFields(logrus.Fields{
		"func": "redis.PushToTail",
		"key":  k,
	}).Infof("pushed: %d", len(posts))

	return nil
}

func GetCachedPosts(actx *appctx.Context, key string) ([]*reddit.Post, error) {
	k := normalizeKey(key)

	log := actx.Logger.WithFields(logrus.Fields{
		"func": "redis.GetCachedPosts",
		"key":  k,
	})
	log.Info("getting cached posts...")

	strings, err := client.LRange(actx.Context, k, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	posts := make([]*reddit.Post, 0)
	for _, s := range strings {
		p := reddit.Post{}
		if err := json.Unmarshal([]byte(s), &p); err != nil {
			return nil, fmt.Errorf("failed to unmarshal post: %v", err)
		}
		posts = append(posts, &p)
	}

	log.Infof("got posts: %d", len(posts))

	return posts, nil
}
