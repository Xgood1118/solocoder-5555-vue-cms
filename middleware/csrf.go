package middleware

import (
	"net/http"
	"sync"

	"cms/utils"

	"github.com/gin-gonic/gin"
)

var (
	csrfTokens   = make(map[string]bool)
	csrfTokensMu sync.Mutex
)

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			token = c.PostForm("csrf_token")
		}

		if token == "" || !validateCSRFToken(token) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "invalid CSRF token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GenerateCSRFToken() string {
	token := utils.GenerateCSRFToken()
	csrfTokensMu.Lock()
	csrfTokens[token] = true
	csrfTokensMu.Unlock()
	return token
}

func validateCSRFToken(token string) bool {
	csrfTokensMu.Lock()
	defer csrfTokensMu.Unlock()

	if csrfTokens[token] {
		delete(csrfTokens, token)
		return true
	}
	return false
}

func CSRFTokenHandler(c *gin.Context) {
	token := GenerateCSRFToken()
	c.JSON(http.StatusOK, gin.H{
		"csrf_token": token,
	})
}
