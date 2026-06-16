package handlers

import (
	"net/http"

	"cms/models"
	"cms/store"
	"cms/utils"

	"github.com/gin-gonic/gin"
)

type CreateCategoryRequest struct {
	Name     string          `json:"name" binding:"required"`
	Slug     string          `json:"slug"`
	ParentID string          `json:"parent_id"`
	Template string          `json:"template"`
	SEO      models.CategorySEO `json:"seo"`
}

type UpdateCategoryRequest struct {
	Name     string          `json:"name"`
	Slug     string          `json:"slug"`
	ParentID string          `json:"parent_id"`
	Template string          `json:"template"`
	SEO      models.CategorySEO `json:"seo"`
}

func GetCategories(c *gin.Context) {
	categories, err := store.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range categories {
		categories[i].Name = utils.SanitizeHTML(categories[i].Name)
		categories[i].SEO.Title = utils.SanitizeHTML(categories[i].SEO.Title)
		categories[i].SEO.Description = utils.SanitizeHTML(categories[i].SEO.Description)
		categories[i].SEO.Keywords = utils.SanitizeHTML(categories[i].SEO.Keywords)
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

func GetCategory(c *gin.Context) {
	id := c.Param("id")

	category, err := store.GetCategoryByID(id)
	if err != nil {
		category, err = store.GetCategoryBySlug(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
	}

	category.Name = utils.SanitizeHTML(category.Name)
	category.SEO.Title = utils.SanitizeHTML(category.SEO.Title)
	category.SEO.Description = utils.SanitizeHTML(category.SEO.Description)
	category.SEO.Keywords = utils.SanitizeHTML(category.SEO.Keywords)

	articles, _ := store.GetArticlesByCategoryPath(category.Path)

	c.JSON(http.StatusOK, gin.H{
		"category": category,
		"articles": articles,
	})
}

func CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &models.Category{
		Name:     req.Name,
		Slug:     req.Slug,
		ParentID: req.ParentID,
		Template: req.Template,
		SEO:      req.SEO,
	}

	if err := store.CreateCategory(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "create", "category", category.ID, "created category: "+category.Name)

	c.JSON(http.StatusCreated, gin.H{"category": category})
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	category, err := store.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.ParentID != "" || req.ParentID == "" && req.Name != "" {
		category.ParentID = req.ParentID
	}
	if req.Template != "" {
		category.Template = req.Template
	}
	if req.SEO.Title != "" || req.SEO.Description != "" || req.SEO.Keywords != "" {
		category.SEO = req.SEO
	}

	if err := store.UpdateCategory(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "update", "category", category.ID, "updated category: "+category.Name)

	c.JSON(http.StatusOK, gin.H{"category": category})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	category, err := store.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	if err := store.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "delete", "category", id, "deleted category: "+category.Name)

	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}

func GetCategoryTree(c *gin.Context) {
	categories, err := store.GetCategoryTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tree := buildCategoryTree(categories, "")

	c.JSON(http.StatusOK, gin.H{"tree": tree})
}

func buildCategoryTree(categories []models.Category, parentID string) []map[string]interface{} {
	var tree []map[string]interface{}

	for _, cat := range categories {
		if cat.ParentID == parentID {
			node := map[string]interface{}{
				"id":       cat.ID,
				"name":     cat.Name,
				"slug":     cat.Slug,
				"path":     cat.Path,
				"template": cat.Template,
				"seo":      cat.SEO,
				"children": buildCategoryTree(categories, cat.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}

func GetCategoryPath(c *gin.Context) {
	id := c.Param("id")

	path, err := store.GetCategoryPath(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"path": path})
}
