package models

import "time"

type ArticleStatus string

const (
	StatusDraft     ArticleStatus = "draft"
	StatusPublished ArticleStatus = "published"
	StatusScheduled ArticleStatus = "scheduled"
)

type ArticlePermission string

const (
	PermPublic ArticlePermission = "public"
	PermLogin  ArticlePermission = "login"
	PermRole   ArticlePermission = "role"
	PermClosed ArticlePermission = "closed"
)

type Article struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Slug        string            `json:"slug"`
	Content     string            `json:"content"`
	Summary     string            `json:"summary"`
	Cover       string            `json:"cover"`
	CategoryID  string            `json:"category_id"`
	Tags        []string          `json:"tags"`
	AuthorID    string            `json:"author_id"`
	Status      ArticleStatus     `json:"status"`
	PublishAt   *time.Time        `json:"publish_at,omitempty"`
	Permissions ArticlePermission `json:"permissions"`
	ViewCount   int               `json:"view_count"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Versions    []ArticleVersion  `json:"versions,omitempty"`
}
