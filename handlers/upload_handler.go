package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cms/config"
	"cms/utils"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" && ext != ".bmp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported image format"})
		return
	}

	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 10MB)"})
		return
	}

	uploadDir := config.AppConfig.UploadDir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	filename := uuid.New().String() + ext
	filepath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	thumbDir := filepath.Join(uploadDir, "thumbs")
	os.MkdirAll(thumbDir, 0755)
	thumbPath := filepath.Join(thumbDir, filename)

	srcImg, err := imaging.Open(filepath)
	if err == nil {
		thumbImg := imaging.Fit(srcImg, 400, 300, imaging.Lanczos)
		imaging.Save(thumbImg, thumbPath)
	}

	CreateAuditLog(c, "upload", "file", filename, "uploaded image: "+file.Filename)

	fileURL := fmt.Sprintf("/uploads/%s", filename)
	thumbURL := fmt.Sprintf("/uploads/thumbs/%s", filename)

	c.JSON(http.StatusOK, gin.H{
		"url":         fileURL,
		"thumbnail":   thumbURL,
		"filename":    file.Filename,
		"size":        file.Size,
		"uploaded_at": time.Now(),
	})
}

func UploadCover(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported image format"})
		return
	}

	uploadDir := filepath.Join(config.AppConfig.UploadDir, "covers")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	filename := uuid.New().String() + ext
	fullPath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	srcImg, err := imaging.Open(fullPath)
	if err == nil {
		coverImg := imaging.Fill(srcImg, 1200, 630, imaging.Center, imaging.Lanczos)
		imaging.Save(coverImg, fullPath)
	}

	CreateAuditLog(c, "upload", "cover", filename, "uploaded cover: "+file.Filename)

	fileURL := fmt.Sprintf("/uploads/covers/%s", filename)

	c.JSON(http.StatusOK, gin.H{
		"url": fileURL,
	})
}

func UploadByURL(c *gin.Context) {
	type URLUploadRequest struct {
		URL string `json:"url" binding:"required"`
	}

	var req URLUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})
		return
	}

	CreateAuditLog(c, "upload", "url", req.URL, "used external image URL")

	c.JSON(http.StatusOK, gin.H{
		"url": req.URL,
		"external": true,
	})
}

func GenerateSlugPreview(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	slug := utils.Slugify(title)

	c.JSON(http.StatusOK, gin.H{
		"slug": slug,
		"preview": slug,
	})
}
