package models

import "time"

type CommentStatus string

const (
	CommentPending  CommentStatus = "pending"
	CommentApproved CommentStatus = "approved"
	CommentRejected CommentStatus = "rejected"
)

type Comment struct {
	ID          string        `json:"id"`
	ArticleID   string        `json:"article_id"`
	ParentID    string        `json:"parent_id"`
	AuthorName  string        `json:"author_name"`
	AuthorEmail string        `json:"author_email"`
	Content     string        `json:"content"`
	Status      CommentStatus `json:"status"`
	IP          string        `json:"ip"`
	UserAgent   string        `json:"user_agent"`
	CreatedAt   time.Time     `json:"created_at"`
	Depth       int           `json:"depth"`
}
