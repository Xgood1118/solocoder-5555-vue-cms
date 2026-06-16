package handlers

import (
	"net/http"
	"time"

	"cms/models"
	"cms/store"
	"cms/utils"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	articles, err := store.GetAllArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	publishedCount := 0
	totalViews := 0
	for _, a := range articles {
		if a.Status == models.StatusPublished {
			publishedCount++
		}
		totalViews += a.ViewCount
	}

	todayUV, _ := store.GetTodayUV()
	todayPV, _ := store.GetTodayPV()

	visitStats, _ := store.GetVisitStatsByDay(7)

	last7Days := make([]map[string]interface{}, 7)
	now := time.Now()
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -(6 - i))
		dateStr := date.Format("2006-01-02")
		stats := visitStats[dateStr]
		if stats == nil {
			last7Days[i] = map[string]interface{}{
				"date": dateStr,
				"pv":   0,
				"uv":   0,
			}
		} else {
			last7Days[i] = map[string]interface{}{
				"date": dateStr,
				"pv":   stats["pv"],
				"uv":   stats["uv"],
			}
		}
	}

	popularArticles, _ := store.GetPopularArticles(10)
	for i := range popularArticles {
		popularArticles[i].Content = ""
	}

	users, _ := store.GetAllUsers()

	pendingComments, _ := store.GetPendingComments()

	c.JSON(http.StatusOK, gin.H{
		"total_articles":  len(articles),
		"published_count": publishedCount,
		"total_views":     totalViews,
		"today_uv":        todayUV,
		"today_pv":        todayPV,
		"last_7_days":     last7Days,
		"popular_articles": popularArticles,
		"total_users":     len(users),
		"pending_comments": len(pendingComments),
	})
}

func GetAuditLogs(c *gin.Context) {
	logs, err := store.GetAllAuditLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range logs {
		logs[i].Action = utils.SanitizeHTML(logs[i].Action)
		logs[i].Resource = utils.SanitizeHTML(logs[i].Resource)
		logs[i].Detail = utils.SanitizeHTML(logs[i].Detail)
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func GetUsers(c *gin.Context) {
	users, err := store.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range users {
		users[i].PasswordHash = ""
		users[i].Username = utils.SanitizeHTML(users[i].Username)
		users[i].Email = utils.SanitizeHTML(users[i].Email)
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func UpdateUserRole(c *gin.Context) {
	id := c.Param("id")

	type RoleUpdateRequest struct {
		Role models.UserRole `json:"role" binding:"required"`
	}

	var req RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := store.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	oldRole := user.Role
	user.Role = req.Role

	if err := store.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "update_role", "user", id, "changed role from "+string(oldRole)+" to "+string(req.Role))

	user.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	user, err := store.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := store.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CreateAuditLog(c, "delete", "user", id, "deleted user: "+user.Username)

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
