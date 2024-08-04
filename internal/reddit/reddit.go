package reddit

import (
	"github.com/sirupsen/logrus"
	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
)

type Post struct {
	ID    string
	URL   string
	Title string
}

const (
	subreddit  = "linuxmemes"
	userAgent  = "linux-memes-bot"
	limitPosts = 100
)

var client *reddit.Client

func Init() error {
	logrus.Info("initializing reddit client...")

	var err error

	client, err = reddit.NewClient(reddit.Credentials{}, reddit.FromEnv, reddit.WithUserAgent(userAgent))
	if err != nil {
		return err
	}

	logrus.Info("reddit client initialized")
	return nil
}

func GetPosts(ctx *appctx.Context) ([]*Post, error) {
	ctx.Logger.WithField("func", "reddit.GetHotPosts").Info("got new request")

	opts := reddit.ListOptions{
		Limit: limitPosts,
	}
	return getPosts(ctx, opts)
}

func GetPostsBefore(ctx *appctx.Context, before string) ([]*Post, error) {
	ctx.Logger.WithFields(logrus.Fields{
		"func":   "reddit.GetHotPostsBefore",
		"before": before,
	}).Info("got new request")

	opts := reddit.ListOptions{
		Limit:  limitPosts,
		Before: before,
	}
	return getPosts(ctx, opts)
}

func GetPostsAfter(ctx *appctx.Context, after string) ([]*Post, error) {
	ctx.Logger.WithFields(logrus.Fields{
		"func":  "reddit.GetHotPostsAfter",
		"after": after,
	}).Info("got new request")

	opts := reddit.ListOptions{
		Limit: limitPosts,
		After: after,
	}
	return getPosts(ctx, opts)
}

func getPosts(ctx *appctx.Context, opts reddit.ListOptions) ([]*Post, error) {
	posts, _, err := client.Subreddit.HotPosts(ctx.Context, subreddit, &opts)
	if err != nil {
		return nil, err
	}

	ctx.Logger.WithField("func", "reddit.getHotPosts").Infof("got posts: %d", len(posts))
	posts = filterPosts(ctx, posts)

	result := make([]*Post, 0)
	for _, p := range posts {
		result = append(result, &Post{
			ID:    p.FullID,
			URL:   p.URL,
			Title: p.Title,
		})
	}

	return result, nil
}
