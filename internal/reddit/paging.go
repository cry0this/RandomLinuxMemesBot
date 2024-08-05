package reddit

import (
	"errors"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/cry0this/RandomLinuxMemesBot/internal/appctx"
)

type PostType string

const (
	HotType PostType = "hot"
	TopType PostType = "top"
)

type Page struct {
	Ctx   *appctx.Context
	Type  PostType
	Opts  interface{}
	Posts []*Post
}

func (p *Page) GetPosts() ([]*Post, error) {
	if len(p.Posts) > 0 {
		return p.Posts, nil
	}

	var posts []*reddit.Post
	var err error

	switch p.Type {
	case HotType:
		opts := p.Opts.(reddit.ListOptions)
		posts, _, err = client.Subreddit.HotPosts(p.Ctx.Context, subreddit, &opts)
	case TopType:
		opts := p.Opts.(reddit.ListPostOptions)
		posts, _, err = client.Subreddit.TopPosts(p.Ctx.Context, subreddit, &opts)
	}

	if err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, errors.New("no posts found")
	}

	p.Posts = convertPosts(posts)
	return p.Posts, nil
}

func (p *Page) NextPage() (*Page, error) {
	posts, err := p.GetPosts()
	if err != nil {
		return nil, err
	}

	after := posts[len(posts)-1].ID
	var opts interface{}

	switch p.Type {
	case HotType:
		opts = reddit.ListOptions{
			Limit: p.Opts.(reddit.ListOptions).Limit,
			After: after,
		}
	case TopType:
		opts = reddit.ListPostOptions{
			Time: p.Opts.(reddit.ListPostOptions).Time,
			ListOptions: reddit.ListOptions{
				Limit: p.Opts.(reddit.ListPostOptions).ListOptions.Limit,
				After: after,
			},
		}
	}
	return &Page{
		Ctx:  p.Ctx,
		Type: p.Type,
		Opts: opts,
	}, nil
}

func (p *Page) PrevPage() (*Page, error) {
	posts, err := p.GetPosts()
	if err != nil {
		return nil, err
	}

	before := posts[0].ID
	var opts interface{}

	switch p.Type {
	case HotType:
		opts = reddit.ListOptions{
			Limit:  p.Opts.(reddit.ListOptions).Limit,
			Before: before,
		}
	case TopType:
		opts = reddit.ListPostOptions{
			Time: p.Opts.(reddit.ListPostOptions).Time,
			ListOptions: reddit.ListOptions{
				Limit:  p.Opts.(reddit.ListPostOptions).ListOptions.Limit,
				Before: before,
			},
		}
	}
	return &Page{
		Ctx:  p.Ctx,
		Type: p.Type,
		Opts: opts,
	}, nil
}

func NewHotPage(actx *appctx.Context) *Page {
	actx.Logger.WithField("func", "reddit.NewHotPage").Info("creating new page")
	opts := reddit.ListOptions{
		Limit: maxLimit,
	}
	return &Page{
		Ctx:  actx,
		Type: HotType,
		Opts: opts,
	}
}

func NewTopPage(actx *appctx.Context, period string) *Page {
	actx.Logger.WithField("func", "reddit.NewTopPage").Info("creating new page")
	opts := reddit.ListPostOptions{
		Time: period,
		ListOptions: reddit.ListOptions{
			Limit: maxLimit,
		},
	}
	return &Page{
		Ctx:  actx,
		Type: TopType,
		Opts: opts,
	}
}
