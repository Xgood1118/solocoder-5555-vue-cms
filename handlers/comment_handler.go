package handlers

import (
	"net/http"

	"cms/models"
	"cms/store"
	"cms/utils"

	"github.com/gin-gonic/gin"
)

type CreateCommentRequest struct {
	ArticleID   string `json:"article_id" binding:"required"`
	ParentID    string `json:"parent_id"`
	AuthorName  string `json:"author_name" binding:"required"`
	AuthorEmail string `json:"author_email" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

func GetComments(c *gin.Context) {
	articleID := c.Query("article_id")
	status := c.Query("status")

	var comments []models.Comment
	var err error

	if articleID != "" {
		comments, err = store.GetApprovedCommentsByArticle(articleID)
	} else if status != "" && isAdmin(c) {
		comments, err = store.GetCommentsByStatus(models.CommentStatus(status))
	} else if isAdmin(c) {
		comments, err = store.GetAllComments()
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range comments {
		comments[i].Content = utils.SanitizeHTML(comments[i].Content)
		comments[i].AuthorName = utils.SanitizeHTML(comments[i].AuthorName)
	}

	nestedComments := buildCommentTree(comments, "")

	c.JSON(http.StatusOK, gin.H{"comments": nestedComments})
}

func buildCommentTree(comments []models.Comment, parentID string) []map[string]interface{} {
	var tree []map[string]interface{}

	for _, comment := range comments {
		if comment.ParentID == parentID {
			node := map[string]interface{}{
				"id":           comment.ID,
				"article_id":   comment.ArticleID,
				"parent_id":    comment.ParentID,
				"author_name":  comment.AuthorName,
				"author_email": comment.AuthorEmail,
				"content":      comment.Content,
				"status":       comment.Status,
				"created_at":   comment.CreatedAt,
				"depth":        comment.Depth,
				"replies":      buildCommentTree(comments, comment.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}

func CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &models.Comment{
		ArticleID:   req.ArticleID,
		ParentID:    req.ParentID,
		AuthorName:  req.AuthorName,
		AuthorEmail: req.AuthorEmail,
		Content:     req.Content,
		IP:          c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
	}

	if err := store.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if comment.Status == models.CommentApproved {
		if comment.ParentID != "" {
			sendReplyNotification(comment)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"comment": comment,
		"message": "comment submitted, awaiting moderation",
	})
}

func ApproveComment(c *gin.Context) {
	id := c.Param("id")

	if err := store.ApproveComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	comment, _ := store.GetCommentByID(id)
	if comment.ParentID != "" {
		sendReplyNotification(comment)
	}

	CreateAuditLog(c, "approve", "comment", id, "approved comment")

	c.JSON(http.StatusOK, gin.H{"message": "comment approved"})
}

func RejectComment(c *gin.Context) {
	id := c.Param("id")

	if err := store.RejectComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "reject", "comment", id, "rejected comment")

	c.JSON(http.StatusOK, gin.H{"message": "comment rejected"})
}

func DeleteComment(c *gin.Context) {
	id := c.Param("id")

	if err := store.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "delete", "comment", id, "deleted comment")

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}

func GetPendingComments(c *gin.Context) {
	comments, err := store.GetPendingComments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range comments {
		comments[i].Content = utils.SanitizeHTML(comments[i].Content)
		comments[i].AuthorName = utils.SanitizeHTML(comments[i].AuthorName)
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func sendReplyNotification(comment *models.Comment) {
}
