package reddit

import (
	"strings"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
)

func filterPosts(ctx *appctx.Context, posts []*reddit.Post) []*reddit.Post {
	filtered := filterNonImagePosts(ctx, posts)
	filtered = filterNSFWPosts(ctx, filtered)
	filtered = filterNonRedditLinks(ctx, filtered)

	ctx.Logger.WithField("func", "reddit.filterPosts").Infof("filtered total: %d", len(posts)-len(filtered))

	return filtered
}

func filterNonImagePosts(ctx *appctx.Context, posts []*reddit.Post) []*reddit.Post {
	filtered := make([]*reddit.Post, 0)

	for _, p := range posts {
		for _, s := range []string{".jpg", ".png", ".gif"} {
			if strings.HasSuffix(p.URL, s) {
				filtered = append(filtered, p)
			}
		}
	}

	ctx.Logger.WithField("func", "reddit.filterNonImagePosts").Infof("filtered: %d", len(posts)-len(filtered))

	return filtered
}

func filterNSFWPosts(ctx *appctx.Context, posts []*reddit.Post) []*reddit.Post {
	filtered := make([]*reddit.Post, 0)

	for _, p := range posts {
		if !p.NSFW {
			filtered = append(filtered, p)
		}
	}

	ctx.Logger.WithField("func", "reddit.filterNSFWPosts").Infof("filtered: %d", len(posts)-len(filtered))

	return filtered
}

func filterNonRedditLinks(ctx *appctx.Context, posts []*reddit.Post) []*reddit.Post {
	filtered := make([]*reddit.Post, 0)

	for _, p := range posts {
		if strings.HasPrefix(p.URL, "https://i.redd.it") {
			filtered = append(filtered, p)
		}
	}

	ctx.Logger.WithField("func", "reddit.filterNonRedditLinks").Infof("filtered: %d", len(posts)-len(filtered))

	return filtered
}
