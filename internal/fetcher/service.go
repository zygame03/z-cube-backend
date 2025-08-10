package fetcher

import (
	"context"
	"strings"
	"z-cube-backend/internal/logger"

	"github.com/mmcdole/gofeed"
	"github.com/panjf2000/ants"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
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

func (s *Service) Run() {
	routes := s.router.FetchableRoutes()
	num := len(routes)
	if num == 0 {
		logger.Info(
			"no fetchable route",
		)
		return
	}

	p, _ := ants.NewPoolWithFunc(10, func(i any) {
		s.FetchRoute(i)
	})
	defer p.Release()

	for _, route := range routes {
		_ = p.Invoke(route)
	}
}

func (s *Service) FetchRoute(route any) {
	ctx := context.Background()
	feed, err := s.parser.ParseURLWithContext(s.router.baseURL, ctx)
	if err != nil {
		logger.Error(
			"parse url failed",
			zap.String("url", s.router.baseURL),
			zap.Error(err),
		)
		return
	}

	if feed == nil {
		logger.Info(
			"feed item failed",
		)
		return
	}

	if len(feed.Items) == 0 {
		logger.Info(
			"no feed items",
		)
		return
	}

	items := make([]*FeedItem, 0, len(feed.Items))
	for _, item := range feed.Items {
		feedItem := s.ItemFormat(item, feed.Link)
		if feedItem != nil {
			items = append(items, feedItem)
		}
	}

	err = s.db.FeedItemsWrite(ctx, items)
	if err != nil {
		logger.Error(
			"repo feed write failed",
			zap.Error(err),
		)
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
