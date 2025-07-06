package infra

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Dbname   string `mapstructure:"dbname" json:"dbname"`
	Port     int    `mapstructure:"port" json:"port"`
	Sslmode  string `mapstructure:"sslmode" json:"sslmode"`
	TimeZone string `mapstructure:"timeZone" json:"timeZone"`
}

func InitDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Dbname, cfg.Port, cfg.Sslmode)

	dialector := postgres.Open(dsn)

	db, err := gorm.Open(dialector)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate()
	if err != nil {
		return nil, err
	}

	return db, nil
}
