package store

import (
	"errors"
	"strings"
	"time"

	"cms/models"
	"cms/utils"

	"github.com/google/uuid"
)

var categoryStore = NewStore("categories.json")

type CategoryList struct {
	Categories []models.Category `json:"categories"`
}

func GetAllCategories() ([]models.Category, error) {
	var list CategoryList
	if err := categoryStore.ReadAll(&list); err != nil {
		return nil, err
	}
	return list.Categories, nil
}

func GetCategoryByID(id string) (*models.Category, error) {
	categories, err := GetAllCategories()
	if err != nil {
		return nil, err
	}
	for _, c := range categories {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, errors.New("category not found")
}

func GetCategoryBySlug(slug string) (*models.Category, error) {
	categories, err := GetAllCategories()
	if err != nil {
		return nil, err
	}
	for _, c := range categories {
		if c.Slug == slug {
			return &c, nil
		}
	}
	return nil, errors.New("category not found")
}

func GetChildCategories(parentID string) ([]models.Category, error) {
	categories, err := GetAllCategories()
	if err != nil {
		return nil, err
	}

	var children []models.Category
	for _, c := range categories {
		if c.ParentID == parentID {
			children = append(children, c)
		}
	}
	return children, nil
}

func GetCategoryTree() ([]models.Category, error) {
	categories, err := GetAllCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func GetCategoryPath(categoryID string) ([]models.Category, error) {
	var path []models.Category
	currentID := categoryID

	for currentID != "" {
		cat, err := GetCategoryByID(currentID)
		if err != nil {
			break
		}
		path = append([]models.Category{*cat}, path...)
		currentID = cat.ParentID
	}

	return path, nil
}

func GetAllDescendantIDs(categoryID string) ([]string, error) {
	var ids []string
	children, err := GetChildCategories(categoryID)
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		ids = append(ids, child.ID)
		descendants, err := GetAllDescendantIDs(child.ID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, descendants...)
	}

	return ids, nil
}

func CreateCategory(category *models.Category) error {
	categories, err := GetAllCategories()
	if err != nil {
		return err
	}

	category.ID = uuid.New().String()

	if category.Slug == "" {
		category.Slug = utils.GenerateUniqueSlug(category.Name, func(slug string) bool {
			_, err := GetCategoryBySlug(slug)
			return err == nil
		})
	}

	if category.ParentID != "" {
		parent, err := GetCategoryByID(category.ParentID)
		if err != nil {
			return errors.New("parent category not found")
		}
		category.Path = parent.Path + "/" + category.Slug
	} else {
		category.Path = "/" + category.Slug
	}

	if category.Template == "" {
		category.Template = "default"
	}

	categories = append(categories, *category)
	return categoryStore.WriteAll(CategoryList{Categories: categories})
}

func UpdateCategory(category *models.Category) error {
	categories, err := GetAllCategories()
	if err != nil {
		return err
	}

	found := false
	for i, c := range categories {
		if c.ID == category.ID {
			if category.Name != c.Name || category.Slug != c.Slug || category.ParentID != c.ParentID {
				if category.Slug == "" {
					category.Slug = utils.GenerateUniqueSlug(category.Name, func(slug string) bool {
						existing, err := GetCategoryBySlug(slug)
						return err == nil && existing.ID != category.ID
					})
				}

				if category.ParentID != "" {
					parent, err := GetCategoryByID(category.ParentID)
					if err != nil {
						return errors.New("parent category not found")
					}
					category.Path = parent.Path + "/" + category.Slug
				} else {
					category.Path = "/" + category.Slug
				}

				updateChildPaths(categories, c.ID, category.Path)
			}

			categories[i] = *category
			found = true
			break
		}
	}

	if !found {
		return errors.New("category not found")
	}

	return categoryStore.WriteAll(CategoryList{Categories: categories})
}

func updateChildPaths(categories []models.Category, parentID string, newParentPath string) {
	for i, c := range categories {
		if c.ParentID == parentID {
			newPath := newParentPath + "/" + c.Slug
			categories[i].Path = newPath
			updateChildPaths(categories, c.ID, newPath)
		}
	}
}

func DeleteCategory(id string) error {
	categories, err := GetAllCategories()
	if err != nil {
		return err
	}

	descendantIDs, err := GetAllDescendantIDs(id)
	if err != nil {
		return err
	}

	deleteIDs := append([]string{id}, descendantIDs...)

	deleteMap := make(map[string]bool)
	for _, did := range deleteIDs {
		deleteMap[did] = true
	}

	var newCategories []models.Category
	for _, c := range categories {
		if !deleteMap[c.ID] {
			newCategories = append(newCategories, c)
		}
	}

	return categoryStore.WriteAll(CategoryList{Categories: newCategories})
}

func GetArticlesByCategoryPath(path string) ([]models.Article, error) {
	categories, err := GetAllCategories()
	if err != nil {
		return nil, err
	}

	var categoryIDs []string
	for _, c := range categories {
		if strings.HasPrefix(c.Path, path) || c.Path == path {
			categoryIDs = append(categoryIDs, c.ID)
		}
	}

	var allArticles []models.Article
	articles, err := GetPublishedArticles()
	if err != nil {
		return nil, err
	}

	for _, a := range articles {
		for _, cid := range categoryIDs {
			if a.CategoryID == cid {
				allArticles = append(allArticles, a)
				break
			}
		}
	}

	return allArticles, nil
}

var articleViewTimes = make(map[string][]time.Time)

func RecordArticleView(articleID string) {
	now := time.Now()
	articleViewTimes[articleID] = append(articleViewTimes[articleID], now)
}

func GetTodayViewCount(articleID string) int {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	count := 0
	for _, t := range articleViewTimes[articleID] {
		if t.After(today) {
			count++
		}
	}
	return count
}

func GetVisitStats(days int) (map[string]int, error) {
	visits, err := GetAllVisitLogs()
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int)
	now := time.Now()

	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		stats[dateStr] = 0
	}

	for _, v := range visits {
		dateStr := v.CreatedAt.Format("2006-01-02")
		if _, exists := stats[dateStr]; exists {
			stats[dateStr]++
		}
	}

	return stats, nil
}
