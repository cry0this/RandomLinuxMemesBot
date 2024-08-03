package reddit

import (
	"strings"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func FilterPosts(posts []*reddit.Post) []*reddit.Post {
	filtered := filterNonImagePosts(posts)
	filtered = filterNSFWPosts(filtered)
	filtered = filterNonRedditLinks(filtered)
	return filtered
}

func filterNonImagePosts(posts []*reddit.Post) []*reddit.Post {
	filtered := make([]*reddit.Post, 0)

	for _, p := range posts {
		for _, s := range []string{".jpg", ".png", ".gif"} {
			if strings.HasSuffix(p.URL, s) {
				filtered = append(filtered, p)
			}
		}
	}

	return filtered
}

func filterNSFWPosts(posts []*reddit.Post) []*reddit.Post {
	filtered := make([]*reddit.Post, 0)

	for _, p := range posts {
		if !p.NSFW {
			filtered = append(filtered, p)
		}
	}

	return filtered
}

func filterNonRedditLinks(posts []*reddit.Post) []*reddit.Post {
	filtered := make([]*reddit.Post, 0)

	for _, p := range posts {
		if strings.HasPrefix(p.URL, "https://i.redd.it") {
			filtered = append(filtered, p)
		}
	}

	return filtered
}
