package models

import "time"

type UserRole string

const (
	RoleGuest  UserRole = "guest"
	RoleAuthor UserRole = "author"
	RoleEditor UserRole = "editor"
	RoleAdmin  UserRole = "admin"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         UserRole  `json:"role"`
	Avatar       string    `json:"avatar"`
	CreatedAt    time.Time `json:"created_at"`
	GitHubID     string    `json:"github_id,omitempty"`
}
