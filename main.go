package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/cry0this/RandomLinuxMemesBot/internal/discord"
	"github.com/cry0this/RandomLinuxMemesBot/internal/memes"
	"github.com/cry0this/RandomLinuxMemesBot/internal/reddit"
	"github.com/cry0this/RandomLinuxMemesBot/internal/redis"
)

var ctx context.Context = context.Background()

func init() {
	logrus.Info("initializing...")

	if err := redis.Init(ctx); err != nil {
		logrus.WithError(err).Fatal("failed to init redis")
	}

	if err := reddit.Init(); err != nil {
		logrus.WithError(err).Fatal("failed to init reddit")
	}

	if err := memes.Init(ctx); err != nil {
		logrus.WithError(err).Fatal("failed to init memes")
	}

	discord.Init(&ctx)

	logrus.Info("app initialized")
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	logrus.Info("closing app...")
	discord.Cleanup()
}
