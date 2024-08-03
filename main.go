package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/cry0this/RandomLinuxMemesBot/internal/discord"
	"github.com/cry0this/RandomLinuxMemesBot/internal/memes"
	"github.com/cry0this/RandomLinuxMemesBot/internal/redis"
	"github.com/sirupsen/logrus"
)

const NSFW = false

var ctx context.Context = context.Background()

func init() {
	logrus.Info("initializing...")

	redis.Init(ctx, 4*time.Hour)
	memes.Init(NSFW)
	discord.Init(&ctx)
}

func main() {
	logrus.Info("application started!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	logrus.Info("cleaning up...")
	discord.Cleanup()
}
