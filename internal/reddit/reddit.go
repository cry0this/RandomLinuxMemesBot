package reddit

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/vartanbeno/go-reddit/v2/reddit"
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

	client, err = reddit.NewReadonlyClient(reddit.WithUserAgent(userAgent))
	if err != nil {
		return err
	}

	logrus.Info("reddit client initialized")
	return nil
}

func GetPosts(ctx context.Context) ([]*Post, error) {
	opts := reddit.ListOptions{
		Limit: limitPosts,
	}
	return getPosts(ctx, opts)
}

func GetPostsBefore(ctx context.Context, before string) ([]*Post, error) {
	opts := reddit.ListOptions{
		Limit:  limitPosts,
		Before: before,
	}
	return getPosts(ctx, opts)
}

func GetPostsAfter(ctx context.Context, after string) ([]*Post, error) {
	opts := reddit.ListOptions{
		Limit: limitPosts,
		After: after,
	}
	return getPosts(ctx, opts)
}

func getPosts(ctx context.Context, opts reddit.ListOptions) ([]*Post, error) {
	posts, _, err := client.Subreddit.HotPosts(ctx, subreddit, &opts)
	if err != nil {
		return nil, err
	}

	posts = FilterPosts(posts)

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
