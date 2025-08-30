package user

import (
	"context"
	"z-cube-backend/internal/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) save(ctx context.Context, user *User) (int, error) {
	err := r.db.WithContext(ctx).
		Save(user).
		Error
	if err != nil {
		logger.Error(
			"save user failed",
			zap.Int("id", user.ID),
			zap.Error(err),
		)
		return 0, err
	}
	logger.Info(
		"save user successfully",
		zap.Int("id", user.ID),
	)
	return user.ID, nil
}

func (r *repo) getByUsername(ctx context.Context, name string) (*User, error) {
	var user User
	err := r.db.
		WithContext(ctx).
		Model(&User{}).
		Where("username = ?", name).
		First(&user).
		Error
	if err != nil {
		logger.Error(
			"get user by username failed",
			zap.String("username", name),
			zap.Error(err),
		)
		return nil, err
	}
	logger.Info(
		"get user by username successfully",
		zap.Int("id", user.ID),
	)

	return &user, nil
}
