package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"cms/models"
	"cms/store"
	"cms/utils"

	"github.com/gin-gonic/gin"
)

type CreateArticleRequest struct {
	Title       string                `json:"title" binding:"required"`
	Content     string                `json:"content" binding:"required"`
	Summary     string                `json:"summary"`
	Cover       string                `json:"cover"`
	CategoryID  string                `json:"category_id"`
	Tags        []string              `json:"tags"`
	Status      models.ArticleStatus  `json:"status"`
	Permissions models.ArticlePermission `json:"permissions"`
	PublishAt   *time.Time            `json:"publish_at"`
}

type UpdateArticleRequest struct {
	Title       string                `json:"title"`
	Content     string                `json:"content"`
	Summary     string                `json:"summary"`
	Cover       string                `json:"cover"`
	CategoryID  string                `json:"category_id"`
	Tags        []string              `json:"tags"`
	Status      models.ArticleStatus  `json:"status"`
	Permissions models.ArticlePermission `json:"permissions"`
	PublishAt   *time.Time            `json:"publish_at"`
	Slug        string                `json:"slug"`
}

func GetArticles(c *gin.Context) {
	status := c.Query("status")
	categoryID := c.Query("category_id")
	tag := c.Query("tag")
	authorID := c.Query("author_id")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var articles []models.Article
	var err error

	if search != "" {
		articles, err = store.SearchArticles(search)
	} else if status != "" {
		if status == "published" {
			articles, err = store.GetPublishedArticles()
		} else if status == "all" && isAdmin(c) {
			articles, err = store.GetAllArticles()
		} else {
			articles, err = store.GetPublishedArticles()
		}
	} else if categoryID != "" {
		articles, err = store.GetArticlesByCategory(categoryID)
	} else if tag != "" {
		articles, err = store.GetArticlesByTag(tag)
	} else if authorID != "" && isAdminOrSameUser(c, authorID) {
		articles, err = store.GetArticlesByAuthor(authorID)
	} else {
		articles, err = store.GetPublishedArticles()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range articles {
		articles[i].Content = utils.SanitizeHTML(articles[i].Content)
		articles[i].Summary = utils.SanitizeHTML(articles[i].Summary)
		articles[i].Title = utils.SanitizeHTML(articles[i].Title)
	}

	total := len(articles)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		articles = []models.Article{}
	} else if end > total {
		articles = articles[start:]
	} else {
		articles = articles[start:end]
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func GetArticle(c *gin.Context) {
	id := c.Param("id")

	article, err := store.GetArticleByID(id)
	if err != nil {
		article, err = store.GetArticleBySlug(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
			return
		}
	}

	if article.Status != models.StatusPublished && !canViewDraft(c, article) {
		c.JSON(http.StatusForbidden, gin.H{"error": "article not published"})
		return
	}

	if article.Status == models.StatusPublished {
		store.IncrementViewCount(id)
		RecordVisit(c, article.ID)
	}

	article.Content = utils.SanitizeHTML(article.Content)
	article.Summary = utils.SanitizeHTML(article.Summary)
	article.Title = utils.SanitizeHTML(article.Title)

	c.JSON(http.StatusOK, gin.H{"article": article})
}

func CreateArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	article := &models.Article{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		Cover:       req.Cover,
		CategoryID:  req.CategoryID,
		Tags:        req.Tags,
		AuthorID:    userID.(string),
		Status:      req.Status,
		Permissions: req.Permissions,
		PublishAt:   req.PublishAt,
	}

	if article.Status == "" {
		article.Status = models.StatusDraft
	}

	if article.Permissions == "" {
		article.Permissions = models.PermPublic
	}

	if err := store.CreateArticle(article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, tagName := range req.Tags {
		store.GetOrCreateTag(tagName)
	}

	version := &models.ArticleVersion{
		Title:     article.Title,
		Content:   article.Content,
		Summary:   article.Summary,
		CreatedBy: userID.(string),
	}
	store.AddArticleVersion(article.ID, version)

	CreateAuditLog(c, "create", "article", article.ID, "created article: "+article.Title)

	c.JSON(http.StatusCreated, gin.H{"article": article})
}

func UpdateArticle(c *gin.Context) {
	id := c.Param("id")

	article, err := store.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	if !canEditArticle(c, article) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Summary != "" {
		article.Summary = req.Summary
	}
	if req.Cover != "" {
		article.Cover = req.Cover
	}
	if req.CategoryID != "" {
		article.CategoryID = req.CategoryID
	}
	if req.Tags != nil {
		article.Tags = req.Tags
	}
	if req.Status != "" {
		article.Status = req.Status
	}
	if req.Permissions != "" {
		article.Permissions = req.Permissions
	}
	if req.PublishAt != nil {
		article.PublishAt = req.PublishAt
	}
	if req.Slug != "" {
		article.Slug = req.Slug
	}

	if err := store.UpdateArticle(article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, tagName := range req.Tags {
		store.GetOrCreateTag(tagName)
	}

	userID, _ := c.Get("user_id")
	version := &models.ArticleVersion{
		Title:     article.Title,
		Content:   article.Content,
		Summary:   article.Summary,
		CreatedBy: userID.(string),
	}
	store.AddArticleVersion(article.ID, version)

	CreateAuditLog(c, "update", "article", article.ID, "updated article: "+article.Title)

	c.JSON(http.StatusOK, gin.H{"article": article})
}

func DeleteArticle(c *gin.Context) {
	id := c.Param("id")

	article, err := store.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	if !canEditArticle(c, article) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	if err := store.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "delete", "article", id, "deleted article: "+article.Title)

	c.JSON(http.StatusOK, gin.H{"message": "article deleted"})
}

func GetArticleVersions(c *gin.Context) {
	id := c.Param("id")

	article, err := store.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	if !canViewDraft(c, article) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	versions, err := store.GetArticleVersions(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

func RevertArticleVersion(c *gin.Context) {
	id := c.Param("id")
	versionID := c.Param("version_id")

	article, err := store.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	if !canEditArticle(c, article) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	if err := store.RevertToVersion(id, versionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "revert", "article", id, "reverted to version: "+versionID)

	updatedArticle, _ := store.GetArticleByID(id)
	c.JSON(http.StatusOK, gin.H{"article": updatedArticle})
}

func SearchArticles(c *gin.Context) {
	keyword := c.Query("q")

	articles, err := store.SearchArticles(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range articles {
		articles[i].Title = utils.HighlightText(utils.SanitizeHTML(articles[i].Title), keyword)
		articles[i].Summary = utils.HighlightText(utils.SanitizeHTML(articles[i].Summary), keyword)
		articles[i].Content = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"keyword":  keyword,
		"total":    len(articles),
	})
}

func isAdmin(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}
	return role.(string) == string(models.RoleAdmin) || role.(string) == string(models.RoleEditor)
}

func isAdminOrSameUser(c *gin.Context, authorID string) bool {
	if isAdmin(c) {
		return true
	}
	userID, exists := c.Get("user_id")
	if !exists {
		return false
	}
	return userID.(string) == authorID
}

func canViewDraft(c *gin.Context, article *models.Article) bool {
	if article.Status == models.StatusPublished {
		return true
	}
	userID, exists := c.Get("user_id")
	if !exists {
		return false
	}
	if isAdmin(c) {
		return true
	}
	return userID.(string) == article.AuthorID
}

func canEditArticle(c *gin.Context, article *models.Article) bool {
	if isAdmin(c) {
		return true
	}
	userID, exists := c.Get("user_id")
	if !exists {
		return false
	}
	role, _ := c.Get("role")
	if role.(string) == string(models.RoleEditor) {
		return true
	}
	return userID.(string) == article.AuthorID
}

func RecordVisit(c *gin.Context, articleID string) {
	log := &models.VisitLog{
		ArticleID: articleID,
		Path:      c.Request.URL.Path,
		IP:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Referer:   c.GetHeader("Referer"),
	}
	store.CreateVisitLog(log)
}

func GenerateSlug(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	slug := utils.GenerateUniqueSlug(title, func(s string) bool {
		_, err := store.GetArticleBySlug(s)
		return err == nil
	})

	c.JSON(http.StatusOK, gin.H{"slug": slug})
}
