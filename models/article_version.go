package models

import "time"

type ArticleVersion struct {
	ID        string    `json:"id"`
	ArticleID string    `json:"article_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Summary   string    `json:"summary"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
}
