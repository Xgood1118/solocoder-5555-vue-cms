package handlers

import (
	"net/http"

	"cms/models"
	"cms/store"
	"cms/utils"

	"github.com/gin-gonic/gin"
)

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

type UpdateTagRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func GetTags(c *gin.Context) {
	tags, err := store.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range tags {
		tags[i].Name = utils.SanitizeHTML(tags[i].Name)
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

func GetTag(c *gin.Context) {
	id := c.Param("id")

	tag, err := store.GetTagByID(id)
	if err != nil {
		tag, err = store.GetTagBySlug(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}
	}

	tag.Name = utils.SanitizeHTML(tag.Name)

	articles, _ := store.GetArticlesByTag(tag.Name)

	c.JSON(http.StatusOK, gin.H{
		"tag":      tag,
		"articles": articles,
	})
}

func CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &models.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := store.CreateTag(tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "create", "tag", tag.ID, "created tag: "+tag.Name)

	c.JSON(http.StatusCreated, gin.H{"tag": tag})
}

func UpdateTag(c *gin.Context) {
	id := c.Param("id")

	tag, err := store.GetTagByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Slug != "" {
		tag.Slug = req.Slug
	}

	if err := store.UpdateTag(tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "update", "tag", tag.ID, "updated tag: "+tag.Name)

	c.JSON(http.StatusOK, gin.H{"tag": tag})
}

func DeleteTag(c *gin.Context) {
	id := c.Param("id")

	tag, err := store.GetTagByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	if err := store.DeleteTag(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "delete", "tag", id, "deleted tag: "+tag.Name)

	c.JSON(http.StatusOK, gin.H{"message": "tag deleted"})
}

func GetTagCloud(c *gin.Context) {
	tagCloud, err := store.GetTagCloud()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": tagCloud})
}
