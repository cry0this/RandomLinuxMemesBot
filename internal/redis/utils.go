package redis

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/cry0this/RandomLinuxMemesBot/internal/reddit"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
)

const (
	commonPrefix = "linux_meme_bot_cache"
	guildPrefix  = "guild"
)

func preparePosts(ctx *appctx.Context, posts []*reddit.Post) ([]string, error) {
	strings := make([]string, 0)
	for _, p := range posts {
		b, err := json.Marshal(p)
		if err != nil {
			ctx.Logger.WithField("func", "redis.preparePosts").WithError(err).Errorf("failed to marshal post: %v", p)
			return nil, err
		}
		strings = append(strings, string(b))
	}

	slices.Reverse(strings)

	return strings, nil
}

func normalizeKey(key string) string {
	if key == "" {
		return commonPrefix
	}

	return fmt.Sprintf("%s:%s:%s", commonPrefix, guildPrefix, key)
}
