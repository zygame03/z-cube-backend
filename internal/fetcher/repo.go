package fetcher

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) FeedItemsWrite(ctx context.Context, items []*FeedItem) error {
	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		if item == nil {
			continue
		}

		err := r.db.
			WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoNothing: true,
			}).Create(item).Error
		if err != nil {
			return err
		}
	}

	return nil
}
