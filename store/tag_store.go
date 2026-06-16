package store

import (
	"errors"

	"cms/models"
	"cms/utils"

	"github.com/google/uuid"
)

var tagStore = NewStore("tags.json")

type TagList struct {
	Tags []models.Tag `json:"tags"`
}

func GetAllTags() ([]models.Tag, error) {
	var list TagList
	if err := tagStore.ReadAll(&list); err != nil {
		return nil, err
	}
	return list.Tags, nil
}

func GetTagByID(id string) (*models.Tag, error) {
	tags, err := GetAllTags()
	if err != nil {
		return nil, err
	}
	for _, t := range tags {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, errors.New("tag not found")
}

func GetTagBySlug(slug string) (*models.Tag, error) {
	tags, err := GetAllTags()
	if err != nil {
		return nil, err
	}
	for _, t := range tags {
		if t.Slug == slug {
			return &t, nil
		}
	}
	return nil, errors.New("tag not found")
}

func GetTagByName(name string) (*models.Tag, error) {
	tags, err := GetAllTags()
	if err != nil {
		return nil, err
	}
	for _, t := range tags {
		if t.Name == name {
			return &t, nil
		}
	}
	return nil, errors.New("tag not found")
}

func CreateTag(tag *models.Tag) error {
	tags, err := GetAllTags()
	if err != nil {
		return err
	}

	existing, _ := GetTagByName(tag.Name)
	if existing != nil {
		return errors.New("tag already exists")
	}

	tag.ID = uuid.New().String()
	tag.Slug = utils.Slugify(tag.Name)

	tags = append(tags, *tag)
	return tagStore.WriteAll(TagList{Tags: tags})
}

func GetOrCreateTag(name string) (*models.Tag, error) {
	tag, err := GetTagByName(name)
	if err == nil {
		return tag, nil
	}

	newTag := &models.Tag{Name: name}
	if err := CreateTag(newTag); err != nil {
		return nil, err
	}

	return GetTagByName(name)
}

func UpdateTag(tag *models.Tag) error {
	tags, err := GetAllTags()
	if err != nil {
		return err
	}

	found := false
	for i, t := range tags {
		if t.ID == tag.ID {
			if tag.Name != t.Name {
				tag.Slug = utils.Slugify(tag.Name)
			}
			tags[i] = *tag
			found = true
			break
		}
	}

	if !found {
		return errors.New("tag not found")
	}

	return tagStore.WriteAll(TagList{Tags: tags})
}

func DeleteTag(id string) error {
	tags, err := GetAllTags()
	if err != nil {
		return err
	}

	var newTags []models.Tag
	for _, t := range tags {
		if t.ID != id {
			newTags = append(newTags, t)
		}
	}

	return tagStore.WriteAll(TagList{Tags: newTags})
}

func GetTagCloud() ([]map[string]interface{}, error) {
	tags, err := GetAllTags()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, tag := range tags {
		articles, err := GetArticlesByTag(tag.Name)
		if err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"id":      tag.ID,
			"name":    tag.Name,
			"slug":    tag.Slug,
			"count":   len(articles),
		})
	}

	return result, nil
}
