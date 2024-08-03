package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/cry0this/RandomLinuxMemesBot/internal/reddit"
)

var client *redis.Client

func Init(ctx context.Context) error {
	logrus.Info("initializing redis...")

	url := os.Getenv("REDIS_URL")
	opts, err := redis.ParseURL(url)
	if err != nil {
		return fmt.Errorf("failed to parse REDIS_URL variable: %v", err)
	}

	log := logrus.WithField("url", url)
	log.Info("connecting to redis...")

	client = redis.NewClient(opts)

	if client.Ping(ctx).String() != "ping: PONG" {
		return errors.New("failed to connect redis, check REDIS_URL variable")
	}

	log.Info("redis initialized")
	return nil
}

func AddToCache(ctx context.Context, post *reddit.Post, key string) error {
	posts := []*reddit.Post{post}
	return PushToTail(ctx, posts, key)
}

func PushToHead(ctx context.Context, posts []*reddit.Post, key string) error {
	p, err := preparePosts(posts)
	if err != nil {
		return err
	}

	k := normalizeKey(key)
	if err := client.LPush(ctx, k, p).Err(); err != nil {
		return err
	}

	return nil
}

func PushToTail(ctx context.Context, posts []*reddit.Post, key string) error {
	p, err := preparePosts(posts)
	if err != nil {
		return err
	}

	k := normalizeKey(key)
	if err := client.RPush(ctx, k, p).Err(); err != nil {
		return err
	}

	return nil
}

func GetCachedPosts(ctx context.Context, key string) ([]*reddit.Post, error) {
	k := normalizeKey(key)

	strings, err := client.LRange(ctx, k, 0, -1).Result()
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

	return posts, nil
}

func Expire(ctx context.Context, key string, ttl time.Duration) error {
	k := normalizeKey(key)
	return client.Expire(ctx, k, ttl).Err()
}
