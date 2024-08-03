package redis

import (
	"context"
	"os"
	"time"

	"github.com/cry0this/RandomLinuxMemesBot/internal/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	Client *redis.Client
	TTL    time.Duration
)

func Init(ctx context.Context, ttl time.Duration) {
	log := logrus.WithField("ttl", ttl)
	log.Info("initializing redis...")

	url := os.Getenv("REDIS_URL")
	opts, err := redis.ParseURL(url)
	if err != nil {
		logrus.WithError(err).Fatal("failed to parse redis URL")
	}

	log = log.WithField("url", url)
	log.Info("connecting to redis...")

	Client = redis.NewClient(opts)

	if Client.Ping(ctx).String() != "ping: PONG" {
		logrus.Fatal("error while connecting to redis DB, check redis URL")
	}

	TTL = ttl

	log.Info("redis initialized")
}

func AddToCache(ctx context.Context, guildID string, url string) error {
	len, err := Client.LLen(ctx, guildID).Result()
	if err != nil {
		return err
	}

	err = Client.LPush(ctx, guildID, url).Err()
	if err != nil {
		return err
	}

	if len > 0 {
		return nil
	}

	err = Client.Expire(ctx, guildID, TTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func IsCached(ctx context.Context, guildID string, url string) (bool, error) {
	r, err := Client.LRange(ctx, guildID, 0, -1).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if utils.StringInSlice(url, r) {
		return true, nil
	}

	return false, nil
}
