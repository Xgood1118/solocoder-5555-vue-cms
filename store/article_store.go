package store

import (
	"errors"
	"sort"
	"strings"
	"time"

	"cms/models"
	"cms/utils"

	"github.com/google/uuid"
)

var articleStore = NewStore("articles.json")

type ArticleList struct {
	Articles []models.Article `json:"articles"`
}

func GetAllArticles() ([]models.Article, error) {
	var list ArticleList
	if err := articleStore.ReadAll(&list); err != nil {
		return nil, err
	}
	return list.Articles, nil
}

func GetPublishedArticles() ([]models.Article, error) {
	articles, err := GetAllArticles()
	if err != nil {
		return nil, err
	}

	var published []models.Article
	now := time.Now()
	for _, a := range articles {
		if a.Status == models.StatusPublished {
			published = append(published, a)
		} else if a.Status == models.StatusScheduled && a.PublishAt != nil && a.PublishAt.Before(now) {
			published = append(published, a)
		}
	}

	sort.Slice(published, func(i, j int) bool {
		if published[i].PublishAt != nil && published[j].PublishAt != nil {
			return published[i].PublishAt.After(*published[j].PublishAt)
		}
		return published[i].CreatedAt.After(published[j].CreatedAt)
	})

	return published, nil
}

func GetArticleByID(id string) (*models.Article, error) {
	articles, err := GetAllArticles()
	if err != nil {
		return nil, err
	}
	for _, a := range articles {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("article not found")
}

func GetArticleBySlug(slug string) (*models.Article, error) {
	articles, err := GetAllArticles()
	if err != nil {
		return nil, err
	}
	for _, a := range articles {
		if a.Slug == slug {
			return &a, nil
		}
	}
	return nil, errors.New("article not found")
}

func GetArticlesByCategory(categoryID string) ([]models.Article, error) {
	articles, err := GetPublishedArticles()
	if err != nil {
		return nil, err
	}

	var result []models.Article
	for _, a := range articles {
		if a.CategoryID == categoryID {
			result = append(result, a)
		}
	}
	return result, nil
}

func GetArticlesByTag(tag string) ([]models.Article, error) {
	articles, err := GetPublishedArticles()
	if err != nil {
		return nil, err
	}

	var result []models.Article
	for _, a := range articles {
		for _, t := range a.Tags {
			if t == tag {
				result = append(result, a)
				break
			}
		}
	}
	return result, nil
}

func GetArticlesByAuthor(authorID string) ([]models.Article, error) {
	articles, err := GetAllArticles()
	if err != nil {
		return nil, err
	}

	var result []models.Article
	for _, a := range articles {
		if a.AuthorID == authorID {
			result = append(result, a)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

func SearchArticles(keyword string) ([]models.Article, error) {
	articles, err := GetPublishedArticles()
	if err != nil {
		return nil, err
	}

	if keyword == "" {
		return articles, nil
	}

	keyword = strings.ToLower(keyword)
	var result []models.Article

	for _, a := range articles {
		titleMatch := strings.Contains(strings.ToLower(a.Title), keyword)
		contentMatch := strings.Contains(strings.ToLower(a.Content), keyword)
		tagMatch := false
		for _, t := range a.Tags {
			if strings.Contains(strings.ToLower(t), keyword) {
				tagMatch = true
				break
			}
		}

		author, err := GetUserByID(a.AuthorID)
		authorMatch := false
		if err == nil {
			authorMatch = strings.Contains(strings.ToLower(author.Username), keyword)
		}

		if titleMatch || contentMatch || tagMatch || authorMatch {
			result = append(result, a)
		}
	}

	return result, nil
}

func CreateArticle(article *models.Article) error {
	articles, err := GetAllArticles()
	if err != nil {
		return err
	}

	article.ID = uuid.New().String()
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	article.ViewCount = 0

	if article.Slug == "" {
		article.Slug = utils.GenerateUniqueSlug(article.Title, func(slug string) bool {
			_, err := GetArticleBySlug(slug)
			return err == nil
		})
	}

	if article.Summary == "" {
		article.Summary = utils.TruncateText(utils.ExtractFirstParagraph(article.Content), 200)
	}

	if article.Status == models.StatusPublished && article.PublishAt == nil {
		now := time.Now()
		article.PublishAt = &now
	}

	articles = append(articles, *article)
	return articleStore.WriteAll(ArticleList{Articles: articles})
}

func UpdateArticle(article *models.Article) error {
	articles, err := GetAllArticles()
	if err != nil {
		return err
	}

	found := false
	for i, a := range articles {
		if a.ID == article.ID {
			article.UpdatedAt = time.Now()

			if article.Slug == "" || article.Slug != a.Slug {
				if article.Title != a.Title {
					article.Slug = utils.GenerateUniqueSlug(article.Title, func(slug string) bool {
						existing, err := GetArticleBySlug(slug)
						return err == nil && existing.ID != article.ID
					})
				}
			}

			if article.Summary == "" {
				article.Summary = utils.TruncateText(utils.ExtractFirstParagraph(article.Content), 200)
			}

			if article.Status == models.StatusPublished && article.PublishAt == nil {
				now := time.Now()
				article.PublishAt = &now
			}

			article.CreatedAt = a.CreatedAt
			article.ViewCount = a.ViewCount
			article.AuthorID = a.AuthorID

			articles[i] = *article
			found = true
			break
		}
	}

	if !found {
		return errors.New("article not found")
	}

	return articleStore.WriteAll(ArticleList{Articles: articles})
}

func DeleteArticle(id string) error {
	articles, err := GetAllArticles()
	if err != nil {
		return err
	}

	var newArticles []models.Article
	for _, a := range articles {
		if a.ID != id {
			newArticles = append(newArticles, a)
		}
	}

	return articleStore.WriteAll(ArticleList{Articles: newArticles})
}

func IncrementViewCount(id string) error {
	articles, err := GetAllArticles()
	if err != nil {
		return err
	}

	found := false
	for i, a := range articles {
		if a.ID == id {
			articles[i].ViewCount++
			found = true
			break
		}
	}

	if !found {
		return errors.New("article not found")
	}

	return articleStore.WriteAll(ArticleList{Articles: articles})
}

func GetPopularArticles(limit int) ([]models.Article, error) {
	articles, err := GetPublishedArticles()
	if err != nil {
		return nil, err
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].ViewCount > articles[j].ViewCount
	})

	if len(articles) > limit {
		return articles[:limit], nil
	}
	return articles, nil
}

func GetScheduledArticles() ([]models.Article, error) {
	articles, err := GetAllArticles()
	if err != nil {
		return nil, err
	}

	var scheduled []models.Article
	now := time.Now()
	for _, a := range articles {
		if a.Status == models.StatusScheduled && a.PublishAt != nil && a.PublishAt.Before(now) {
			scheduled = append(scheduled, a)
		}
	}
	return scheduled, nil
}

func AddArticleVersion(articleID string, version *models.ArticleVersion) error {
	articles, err := GetAllArticles()
	if err != nil {
		return err
	}

	version.ID = uuid.New().String()
	version.ArticleID = articleID
	version.CreatedAt = time.Now()

	found := false
	for i, a := range articles {
		if a.ID == articleID {
			version.Version = len(a.Versions) + 1
			articles[i].Versions = append(a.Versions, *version)
			if len(articles[i].Versions) > 10 {
				articles[i].Versions = articles[i].Versions[len(articles[i].Versions)-10:]
			}
			found = true
			break
		}
	}

	if !found {
		return errors.New("article not found")
	}

	return articleStore.WriteAll(ArticleList{Articles: articles})
}

func GetArticleVersions(articleID string) ([]models.ArticleVersion, error) {
	article, err := GetArticleByID(articleID)
	if err != nil {
		return nil, err
	}
	return article.Versions, nil
}

func RevertToVersion(articleID string, versionID string) error {
	articles, err := GetAllArticles()
	if err != nil {
		return err
	}

	found := false
	for i, a := range articles {
		if a.ID == articleID {
			for _, v := range a.Versions {
				if v.ID == versionID {
					articles[i].Title = v.Title
					articles[i].Content = v.Content
					articles[i].Summary = v.Summary
					articles[i].UpdatedAt = time.Now()
					found = true
					break
				}
			}
			break
		}
	}

	if !found {
		return errors.New("version not found")
	}

	return articleStore.WriteAll(ArticleList{Articles: articles})
}
