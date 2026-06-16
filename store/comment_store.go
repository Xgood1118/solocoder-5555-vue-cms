package store

import (
	"errors"
	"sort"
	"time"

	"cms/models"
	"cms/utils"

	"github.com/google/uuid"
)

var commentStore = NewStore("comments.json")

var SensitiveWords = []string{
	"垃圾评论",
	"广告",
	"赌博",
	"色情",
	"viagra",
	"casino",
}

const LinkThreshold = 3

type CommentList struct {
	Comments []models.Comment `json:"comments"`
}

func GetAllComments() ([]models.Comment, error) {
	var list CommentList
	if err := commentStore.ReadAll(&list); err != nil {
		return nil, err
	}
	return list.Comments, nil
}

func GetCommentByID(id string) (*models.Comment, error) {
	comments, err := GetAllComments()
	if err != nil {
		return nil, err
	}
	for _, c := range comments {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, errors.New("comment not found")
}

func GetApprovedCommentsByArticle(articleID string) ([]models.Comment, error) {
	comments, err := GetAllComments()
	if err != nil {
		return nil, err
	}

	var result []models.Comment
	for _, c := range comments {
		if c.ArticleID == articleID && c.Status == models.CommentApproved {
			result = append(result, c)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result, nil
}

func GetPendingComments() ([]models.Comment, error) {
	comments, err := GetAllComments()
	if err != nil {
		return nil, err
	}

	var result []models.Comment
	for _, c := range comments {
		if c.Status == models.CommentPending {
			result = append(result, c)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

func GetCommentsByStatus(status models.CommentStatus) ([]models.Comment, error) {
	comments, err := GetAllComments()
	if err != nil {
		return nil, err
	}

	var result []models.Comment
	for _, c := range comments {
		if c.Status == status {
			result = append(result, c)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

func IsSpam(content string) bool {
	if utils.ContainsSensitiveWords(content, SensitiveWords) {
		return true
	}

	if utils.CountLinks(content) > LinkThreshold {
		return true
	}

	return false
}

func CreateComment(comment *models.Comment) error {
	if !utils.IsValidEmail(comment.AuthorEmail) {
		return errors.New("invalid email")
	}

	if comment.Content == "" {
		return errors.New("content is required")
	}

	comments, err := GetAllComments()
	if err != nil {
		return err
	}

	comment.ID = uuid.New().String()
	comment.CreatedAt = time.Now()

	if comment.ParentID != "" {
		parent, err := GetCommentByID(comment.ParentID)
		if err != nil {
			return errors.New("parent comment not found")
		}
		comment.Depth = parent.Depth + 1
		if comment.Depth > 3 {
			return errors.New("maximum comment depth exceeded (max 3 levels)")
		}
	} else {
		comment.Depth = 1
	}

	if IsSpam(comment.Content) {
		comment.Status = models.CommentRejected
	} else {
		comment.Status = models.CommentPending
	}

	comments = append(comments, *comment)
	return commentStore.WriteAll(CommentList{Comments: comments})
}

func ApproveComment(id string) error {
	comments, err := GetAllComments()
	if err != nil {
		return err
	}

	found := false
	for i, c := range comments {
		if c.ID == id {
			comments[i].Status = models.CommentApproved
			found = true
			break
		}
	}

	if !found {
		return errors.New("comment not found")
	}

	return commentStore.WriteAll(CommentList{Comments: comments})
}

func RejectComment(id string) error {
	comments, err := GetAllComments()
	if err != nil {
		return err
	}

	found := false
	for i, c := range comments {
		if c.ID == id {
			comments[i].Status = models.CommentRejected
			found = true
			break
		}
	}

	if !found {
		return errors.New("comment not found")
	}

	return commentStore.WriteAll(CommentList{Comments: comments})
}

func DeleteComment(id string) error {
	comments, err := GetAllComments()
	if err != nil {
		return err
	}

	childIDs := getChildCommentIDs(comments, id)
	deleteIDs := append([]string{id}, childIDs...)

	deleteMap := make(map[string]bool)
	for _, did := range deleteIDs {
		deleteMap[did] = true
	}

	var newComments []models.Comment
	for _, c := range comments {
		if !deleteMap[c.ID] {
			newComments = append(newComments, c)
		}
	}

	return commentStore.WriteAll(CommentList{Comments: newComments})
}

func getChildCommentIDs(comments []models.Comment, parentID string) []string {
	var ids []string
	for _, c := range comments {
		if c.ParentID == parentID {
			ids = append(ids, c.ID)
			ids = append(ids, getChildCommentIDs(comments, c.ID)...)
		}
	}
	return ids
}

func GetCommentCountByArticle(articleID string) (int, error) {
	comments, err := GetApprovedCommentsByArticle(articleID)
	if err != nil {
		return 0, err
	}
	return len(comments), nil
}

func GetCommentReplies(parentID string) ([]models.Comment, error) {
	comments, err := GetAllComments()
	if err != nil {
		return nil, err
	}

	var replies []models.Comment
	for _, c := range comments {
		if c.ParentID == parentID && c.Status == models.CommentApproved {
			replies = append(replies, c)
		}
	}

	sort.Slice(replies, func(i, j int) bool {
		return replies[i].CreatedAt.Before(replies[j].CreatedAt)
	})

	return replies, nil
}
