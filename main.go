package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/cry0this/RandomLinuxMemesBot/internal/discord"
	"github.com/cry0this/RandomLinuxMemesBot/internal/reddit"
	"github.com/cry0this/RandomLinuxMemesBot/internal/redis"
)

var ctx context.Context = context.Background()

func init() {
	logrus.Info("initializing app...")

	if err := redis.Init(ctx); err != nil {
		logrus.WithError(err).Fatal("failed to init redis")
	}

	if err := reddit.Init(); err != nil {
		logrus.WithError(err).Fatal("failed to init reddit")
	}

	if err := discord.Init(ctx); err != nil {
		logrus.WithError(err).Fatal("failed to init discord")
	}

	logrus.Info("app initialized")
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-sig

	logrus.Info("closing app...")

	if err := discord.Cleanup(); err != nil {
		logrus.WithError(err).Fatal("failed to cleanup discord")
	}
}
