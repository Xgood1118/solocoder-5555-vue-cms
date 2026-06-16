package models

import "time"

type VisitLog struct {
	ID        string    `json:"id"`
	ArticleID string    `json:"article_id,omitempty"`
	Path      string    `json:"path"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
	CreatedAt time.Time `json:"created_at"`
}
