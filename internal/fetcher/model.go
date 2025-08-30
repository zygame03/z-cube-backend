package fetcher

import (
	"crypto/sha256"
	"encoding/hex"
)

type FeedItem struct {
	ID          string `gorm:"primaryKey;size:64"`
	Title       string `gorm:"size:512"`
	Link        string `gorm:"size:1024;index"`
	Description string `gorm:"type:text"`
	Published   string `gorm:"size:64"`
	Author      string `gorm:"size:255"`
	Category    string `gorm:"size:255"`
	Tags        string `gorm:"size:255"`
	Source      string `gorm:"size:512"`
}

func (f *FeedItem) GetId() string {
	h := sha256.New()
	h.Write([]byte(f.Source))
	h.Write([]byte("|"))
	h.Write([]byte(f.Link))
	return hex.EncodeToString(h.Sum(nil))
}
