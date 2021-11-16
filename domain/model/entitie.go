package model

import (
	"time"
)

type ShortUrl struct {
	Id         string    `json:"id" firestore:"id"`
	Url        string    `json:"url" firestore:"url"`
	CreateTime time.Time `json:"createTime,omitempty" firestore:"createTime,omitempty"`
	Enable     bool      `json:"enable" firestore:"enable"`
	Clicks     int64     `json:"clicks" firestore:"clicks"`
}
