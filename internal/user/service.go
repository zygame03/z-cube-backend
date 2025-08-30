package user

import (
	"context"
	"errors"
	"time"
	"z-cube-backend/internal/logger"
	"z-cube-backend/internal/middleware"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound     = errors.New("未知用户")
	ErrUserAlreadyExist = errors.New("用户已存在")
	ErrInvalidPassword  = errors.New("密码错误")
)

type service struct {
	repo *repo
	cfg  func() *Config
}

func NewService(repo *repo, cfg func() *Config) *service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

// 用户注册
// 返回id和错误
func (s *service) register(ctx context.Context, username, password string) (int, error) {
	// 判断用户是否存在
	user, err := s.repo.getByUsername(ctx, username)
	if err != nil {
		logger.Error(
			"get user failed",
			zap.Int("id", user.ID),
			zap.Error(err),
		)
		return 0, err
	}
	if user != nil {
		return user.ID, ErrUserAlreadyExist
	}

	// 密码明文哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// 构造用户并落库
	newUser := User{
		Username: username,
		Password: string(hashedPassword),
	}
	id, err := s.repo.save(ctx, &newUser)
	if err != nil {
		logger.Error(
			"creat user failed",
			zap.Error(err),
		)
		return 0, err
	}
	logger.Info(
		"register successfully",
		zap.Int("id", newUser.ID),
	)

	return id, nil
}

// 用户登录
func (s *service) login(ctx context.Context, username, password string) (string, error) {
	// 用户名查找用户
	user, err := s.repo.getByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", ErrUserNotFound
	}

	// 比对密码哈希值
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}

	// 获取ttl，判空+保底
	ttl := 24 * time.Hour
	cfg := s.cfg()
	if cfg != nil && cfg.UserTokenTTL > 0 {
		ttl = s.cfg().UserTokenTTL
	}

	// 生成token
	token, err := middleware.GenerateToken(user.ID, username, ttl)
	if err != nil {
		logger.Error(
			"generate token failed",
			zap.Int("id", user.ID),
			zap.Error(err),
		)
		return "", err
	}

	return token, nil
}
