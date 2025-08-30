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
	cfg    func() *Config
	parser *gofeed.Parser
	router *Router
}

func NewService(db *gorm.DB, rdb *redis.Client, cfg func() *Config) *Service {
	srv := &Service{
		db:     NewRepo(db),
		rdb:    NewCache(rdb, cfg),
		cfg:    cfg,
		parser: gofeed.NewParser(),
	}
	srv.router = NewRouter(srv.cfg().BaseURL)
	return srv
}

type routeView struct {
	name string
	path string
}

// 注册定时任务
// 感觉需要调整，能实时控制开启和关闭
func (s *Service) RegisterCron(cron *cron.Cron) error {
	_, err := cron.AddFunc("@every "+s.cfg().Interval.String(), func() {
		s.Run()
	})
	if err != nil {
		logger.Error(
			"add cron failed",
			zap.String("func", "Run"),
		)
	}
	logger.Info(
		"add cron",
		zap.String("func", "Run"),
	)
	return nil
}

func (s *Service) Run() {
	// 获取开启的路由
	routes := s.router.FetchableRoutes()
	num := len(routes)
	if num == 0 {
		logger.Info(
			"no fetchable route",
		)
		return
	}

	// 使用协程池并发获取
	p, _ := ants.NewPoolWithFunc(s.cfg().Concurrency, func(i any) {
		s.FetchRoute(i)
	})
	defer p.Release()

	// 提交参数
	for _, route := range routes {
		_ = p.Invoke(route)
	}
}

// 获取路由
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

	// 也许应该先写一份缓存
	// TODO

	// 落库
	// 是否改为异步进行？
	err = s.db.FeedItemsWrite(ctx, items)
	if err != nil {
		logger.Error(
			"repo feed write failed",
			zap.Error(err),
		)
	}
	logger.Info(
		"fetch finished",
		zap.String("name", route.(*routeView).name),
		zap.String("url", route.(*routeView).path),
	)
}

// 数据格式化
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
