package main

import (
	"fmt"
	"log"
	"time"

	"cms/config"
	"cms/handlers"
	"cms/middleware"
	"cms/models"
	"cms/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	config.LoadConfig()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.RateLimitMiddleware(60, time.Minute))

	r.Static("/uploads", "./uploads")
	r.Static("/frontend", "./frontend/dist")

	setupRoutes(r)
	initDefaultData()
	startCronJobs()

	log.Printf("Server starting on port %d...", config.AppConfig.Port)
	log.Printf("Data directory: %s", config.AppConfig.DataDir)
	log.Printf("Upload directory: %s", config.AppConfig.UploadDir)

	r.Run(fmt.Sprintf(":%d", config.AppConfig.Port))
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/csrf-token", handlers.GetCSRFToken)

	auth := api.Group("/auth")
	{
		auth.POST("/register", middleware.CSRFMiddleware(), handlers.Register)
		auth.POST("/login", middleware.CSRFMiddleware(), handlers.Login)
		auth.POST("/refresh", handlers.RefreshToken)
		auth.POST("/logout", middleware.AuthMiddleware(), handlers.Logout)
		auth.GET("/me", middleware.AuthMiddleware(), handlers.GetCurrentUser)
		auth.GET("/github", handlers.GitHubLogin)
		auth.GET("/github/callback", handlers.GitHubCallback)
	}

	articles := api.Group("/articles")
	{
		articles.GET("", middleware.OptionalAuthMiddleware(), handlers.GetArticles)
		articles.GET("/search", handlers.SearchArticles)
		articles.GET("/slug", handlers.GenerateSlug)
		articles.GET("/:id", middleware.OptionalAuthMiddleware(), handlers.GetArticle)
		articles.GET("/:id/versions", middleware.AuthMiddleware(), handlers.GetArticleVersions)
		articles.POST("", middleware.AuthMiddleware(), middleware.CSRFMiddleware(), handlers.CreateArticle)
		articles.PUT("/:id", middleware.AuthMiddleware(), middleware.CSRFMiddleware(), handlers.UpdateArticle)
		articles.DELETE("/:id", middleware.AuthMiddleware(), middleware.CSRFMiddleware(), handlers.DeleteArticle)
		articles.POST("/:id/versions/:version_id/revert", middleware.AuthMiddleware(), middleware.CSRFMiddleware(), handlers.RevertArticleVersion)
	}

	categories := api.Group("/categories")
	{
		categories.GET("", handlers.GetCategories)
		categories.GET("/tree", handlers.GetCategoryTree)
		categories.GET("/:id", handlers.GetCategory)
		categories.GET("/:id/path", handlers.GetCategoryPath)
		categories.POST("", middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleEditor, models.RoleAdmin), middleware.CSRFMiddleware(), handlers.CreateCategory)
		categories.PUT("/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleEditor, models.RoleAdmin), middleware.CSRFMiddleware(), handlers.UpdateCategory)
		categories.DELETE("/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleAdmin), middleware.CSRFMiddleware(), handlers.DeleteCategory)
	}

	tags := api.Group("/tags")
	{
		tags.GET("", handlers.GetTags)
		tags.GET("/cloud", handlers.GetTagCloud)
		tags.GET("/:id", handlers.GetTag)
		tags.POST("", middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleEditor, models.RoleAdmin), middleware.CSRFMiddleware(), handlers.CreateTag)
		tags.PUT("/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleEditor, models.RoleAdmin), middleware.CSRFMiddleware(), handlers.UpdateTag)
		tags.DELETE("/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleAdmin), middleware.CSRFMiddleware(), handlers.DeleteTag)
	}

	comments := api.Group("/comments")
	{
		comments.GET("", middleware.OptionalAuthMiddleware(), handlers.GetComments)
		comments.POST("", middleware.CSRFMiddleware(), handlers.CreateComment)
		comments.GET("/pending", middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleEditor, models.RoleAdmin), handlers.GetPendingComments)
		comments.POST("/:id/approve", middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleEditor, models.RoleAdmin), middleware.CSRFMiddleware(), handlers.ApproveComment)
		comments.POST("/:id/reject", middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleEditor, models.RoleAdmin), middleware.CSRFMiddleware(), handlers.RejectComment)
		comments.DELETE("/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleAdmin), middleware.CSRFMiddleware(), handlers.DeleteComment)
	}

	upload := api.Group("/upload")
	upload.Use(middleware.AuthMiddleware())
	{
		upload.POST("/image", middleware.CSRFMiddleware(), handlers.UploadImage)
		upload.POST("/cover", middleware.CSRFMiddleware(), handlers.UploadCover)
		upload.POST("/url", middleware.CSRFMiddleware(), handlers.UploadByURL)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleEditor, models.RoleAdmin))
	{
		admin.GET("/dashboard", handlers.GetDashboardStats)
		admin.GET("/audit-logs", middleware.RoleMiddleware(models.RoleAdmin), handlers.GetAuditLogs)
		admin.GET("/users", middleware.RoleMiddleware(models.RoleAdmin), handlers.GetUsers)
		admin.PUT("/users/:id/role", middleware.RoleMiddleware(models.RoleAdmin), middleware.CSRFMiddleware(), handlers.UpdateUserRole)
		admin.DELETE("/users/:id", middleware.RoleMiddleware(models.RoleAdmin), middleware.CSRFMiddleware(), handlers.DeleteUser)
	}

	api.GET("/rss", handlers.GetRSSFeed)
	api.GET("/sitemap.xml", handlers.GetSitemap)
	api.POST("/markdown/render", handlers.MarkdownRender)

	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}

func initDefaultData() {
	users, _ := store.GetAllUsers()
	if len(users) == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := &models.User{
			Username:     "admin",
			Email:        "admin@example.com",
			PasswordHash: string(hashedPassword),
			Role:         models.RoleAdmin,
			Avatar:       "",
		}
		store.CreateUser(admin)
		log.Println("Created default admin user: admin / admin123")
	}

	categories, _ := store.GetAllCategories()
	if len(categories) == 0 {
		defaultCat := &models.Category{
			Name:     "默认分类",
			Slug:     "default",
			ParentID: "",
			Path:     "/default",
			Template: "default",
			SEO: models.CategorySEO{
				Title:       "默认分类 - Vue CMS",
				Description: "默认文章分类",
				Keywords:    "博客,文章",
			},
		}
		store.CreateCategory(defaultCat)
		log.Println("Created default category")
	}
}

func startCronJobs() {
	c := cron.New()

	c.AddFunc("*/5 * * * *", func() {
		published := handlers.CheckScheduled()
		if published > 0 {
			log.Printf("Published %d scheduled articles", published)
		}
	})

	c.Start()
	log.Println("Cron jobs started")
}
