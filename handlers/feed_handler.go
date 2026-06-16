package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"cms/models"
	"cms/store"

	"github.com/gin-gonic/gin"
)

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

func GetRSSFeed(c *gin.Context) {
	articles, err := store.GetPublishedArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	baseURL := fmt.Sprintf("%s://%s", c.Request.URL.Scheme, c.Request.Host)

	var items []RSSItem
	for _, article := range articles {
		pubDate := article.CreatedAt.Format(time.RFC1123Z)
		if article.PublishAt != nil {
			pubDate = article.PublishAt.Format(time.RFC1123Z)
		}

		items = append(items, RSSItem{
			Title:       article.Title,
			Link:        fmt.Sprintf("%s/article/%s", baseURL, article.Slug),
			Description: article.Summary,
			PubDate:     pubDate,
			Guid:        article.ID,
		})

		if len(items) >= 50 {
			break
		}
	}

	feed := RSSFeed{
		Version: "2.0",
		Channel: RSSChannel{
			Title:       "Vue CMS Blog",
			Link:        baseURL,
			Description: "Latest articles from the blog",
			Items:       items,
		},
	}

	c.Header("Content-Type", "application/rss+xml; charset=utf-8")
	xmlBytes, _ := xml.MarshalIndent(feed, "", "  ")
	c.String(http.StatusOK, xml.Header+string(xmlBytes))
}

type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

type Sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

func GetSitemap(c *gin.Context) {
	articles, err := store.GetPublishedArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categories, _ := store.GetAllCategories()
	tags, _ := store.GetAllTags()

	baseURL := fmt.Sprintf("%s://%s", c.Request.URL.Scheme, c.Request.Host)

	var urls []SitemapURL

	urls = append(urls, SitemapURL{
		Loc:        baseURL + "/",
		ChangeFreq: "daily",
		Priority:   "1.0",
	})

	for _, cat := range categories {
		urls = append(urls, SitemapURL{
			Loc:        fmt.Sprintf("%s/category/%s", baseURL, cat.Slug),
			ChangeFreq: "weekly",
			Priority:   "0.8",
		})
	}

	for _, tag := range tags {
		urls = append(urls, SitemapURL{
			Loc:        fmt.Sprintf("%s/tag/%s", baseURL, tag.Slug),
			ChangeFreq: "weekly",
			Priority:   "0.6",
		})
	}

	for _, article := range articles {
		lastMod := article.UpdatedAt.Format("2006-01-02")
		urls = append(urls, SitemapURL{
			Loc:        fmt.Sprintf("%s/article/%s", baseURL, article.Slug),
			LastMod:    lastMod,
			ChangeFreq: "monthly",
			Priority:   "0.9",
		})
	}

	sitemap := Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	c.Header("Content-Type", "application/xml; charset=utf-8")
	xmlBytes, _ := xml.MarshalIndent(sitemap, "", "  ")
	c.String(http.StatusOK, xml.Header+string(xmlBytes))
}

func MarkdownRender(c *gin.Context) {
	type RenderRequest struct {
		Content string `json:"content"`
	}

	var req RenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	html := renderMarkdown(req.Content)

	c.JSON(http.StatusOK, gin.H{"html": html})
}

func renderMarkdown(content string) string {
	return content
}

func CheckScheduled() int {
	scheduled, _ := store.GetScheduledArticles()
	published := 0
	for _, article := range scheduled {
		article.Status = models.StatusPublished
		store.UpdateArticle(&article)
		published++
	}
	return published
}
