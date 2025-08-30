package user

import "time"

type Config struct {
	UserTokenTTL time.Duration `json:"user_token_ttl"`
}
