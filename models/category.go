package models

type CategorySEO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
}

type Category struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Slug     string      `json:"slug"`
	ParentID string      `json:"parent_id"`
	Path     string      `json:"path"`
	Template string      `json:"template"`
	SEO      CategorySEO `json:"seo"`
}
