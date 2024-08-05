package reddit

import (
	"github.com/sirupsen/logrus"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Post struct {
	ID    string
	URL   string
	Title string
	NSFW  bool
}

const (
	subreddit = "linuxmemes"
	userAgent = "linux-memes-bot"
	maxLimit  = 100
)

var client *reddit.Client

func Init() error {
	log := logrus.WithField("module", "reddit")
	log.Info("initializing...")

	var err error
	client, err = reddit.NewClient(reddit.Credentials{}, reddit.FromEnv, reddit.WithUserAgent(userAgent))
	if err != nil {
		return err
	}

	log.Info("initialized")
	return nil
}
