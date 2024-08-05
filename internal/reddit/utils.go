package reddit

import (
	"strings"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
)

func FilterPosts(actx *appctx.Context, posts []*Post) []*Post {
	log := actx.Logger.WithField("func", "reddit.FilterPosts")
	log.Infof("got posts: %d", len(posts))

	filtered := filterNonImagePosts(actx, posts)
	filtered = filterNSFWPosts(actx, filtered)
	filtered = filterNonRedditLinks(actx, filtered)

	log.Infof("filtered total: %d", len(posts)-len(filtered))

	return filtered
}

func filterNonImagePosts(actx *appctx.Context, posts []*Post) []*Post {
	filtered := make([]*Post, 0)

	for _, p := range posts {
		for _, s := range []string{".jpg", ".jpeg", ".png", ".gif"} {
			if strings.HasSuffix(p.URL, s) {
				filtered = append(filtered, p)
			}
		}
	}

	actx.Logger.WithField("func", "reddit.filterNonImagePosts").Infof("filtered: %d", len(posts)-len(filtered))

	return filtered
}

func filterNSFWPosts(actx *appctx.Context, posts []*Post) []*Post {
	filtered := make([]*Post, 0)

	for _, p := range posts {
		if !p.NSFW {
			filtered = append(filtered, p)
		}
	}

	actx.Logger.WithField("func", "reddit.filterNSFWPosts").Infof("filtered: %d", len(posts)-len(filtered))

	return filtered
}

func filterNonRedditLinks(actx *appctx.Context, posts []*Post) []*Post {
	filtered := make([]*Post, 0)

	for _, p := range posts {
		if strings.HasPrefix(p.URL, "https://i.redd.it") {
			filtered = append(filtered, p)
		}
	}

	actx.Logger.WithField("func", "reddit.filterNonRedditLinks").Infof("filtered: %d", len(posts)-len(filtered))

	return filtered
}

func convertPosts(posts []*reddit.Post) []*Post {
	result := make([]*Post, 0)
	for _, p := range posts {
		result = append(result, &Post{
			ID:    p.FullID,
			URL:   p.URL,
			Title: p.Title,
			NSFW:  p.NSFW,
		})
	}

	return result
}
