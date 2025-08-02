package fetcher

import (
	"context"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/panjf2000/ants"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Service struct {
	db     *repo
	rdb    *cache
	parser *gofeed.Parser
	router *Router
}

func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{
		db:     NewRepo(db),
		rdb:    NewCache(rdb),
		parser: gofeed.NewParser(),
		router: NewRouter("132.232.238.184:1200"),
	}
}

type routeView struct {
	name string
	path string
}

func (s *Service) RegisterCron(cron *cron.Cron) error {
	_, err := cron.AddFunc("@every 3h", func() {
		s.Run()
	})
	if err != nil {
		return nil
	}
	return nil
}

// 执行 RSS 抓取任务
func (s *Service) Run() {
	// 获取开启的路由
	routes := s.router.FetchableRoutes()
	num := len(routes)
	if num == 0 {
		return
	}

	p, _ := ants.NewPoolWithFunc(10, func(i any) {
		s.FetchRoute(i)
	})
	defer p.Release()

	// 提交任务
	for _, route := range routes {
		_ = p.Invoke(route)
	}
}

// 抓取单个路由的 RSS Feed
func (s *Service) FetchRoute(route any) {
	ctx := context.Background()
	feed, err := s.parser.ParseURLWithContext(s.router.baseURL, ctx)
	if err != nil {
		return
	}

	if feed == nil {
		return
	}

	if len(feed.Items) == 0 {
		return
	}

	items := make([]*FeedItem, 0, len(feed.Items))
	for _, item := range feed.Items {
		feedItem := s.ItemFormat(item, feed.Link)
		if feedItem != nil {
			items = append(items, feedItem)
		}
	}

	// 落库
	err = s.db.FeedItemsWrite(ctx, items)
	if err != nil {
	}
}

// 格式化
func (s *Service) ItemFormat(item *gofeed.Item, source string) *FeedItem {
	if item == nil {
		return nil
	}

	author := ""
	if item.Author != nil {
		author = item.Author.Name
	}

	category := ""
	if len(item.Categories) > 0 {
		category = strings.Join(item.Categories, ",")
	}

	feedItem := &FeedItem{
		Title:       item.Title,
		Link:        item.Link,
		Description: item.Description,
		Published:   item.Published,
		Author:      author,
		Category:    category,
		Source:      source,
	}

	feedItem.ID = feedItem.GetId()

	return feedItem
}
