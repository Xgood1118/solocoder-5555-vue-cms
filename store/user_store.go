package store

import (
	"errors"
	"time"

	"cms/models"

	"github.com/google/uuid"
)

var userStore = NewStore("users.json")

type UserList struct {
	Users []models.User `json:"users"`
}

func GetAllUsers() ([]models.User, error) {
	var list UserList
	if err := userStore.ReadAll(&list); err != nil {
		return nil, err
	}
	return list.Users, nil
}

func GetUserByID(id string) (*models.User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func GetUserByUsername(username string) (*models.User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func GetUserByEmail(email string) (*models.User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func GetUserByGitHubID(githubID string) (*models.User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.GitHubID == githubID {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func CreateUser(user *models.User) error {
	users, err := GetAllUsers()
	if err != nil {
		return err
	}

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	users = append(users, *user)

	return userStore.WriteAll(UserList{Users: users})
}

func UpdateUser(user *models.User) error {
	users, err := GetAllUsers()
	if err != nil {
		return err
	}

	found := false
	for i, u := range users {
		if u.ID == user.ID {
			users[i] = *user
			found = true
			break
		}
	}

	if !found {
		return errors.New("user not found")
	}

	return userStore.WriteAll(UserList{Users: users})
}

func DeleteUser(id string) error {
	users, err := GetAllUsers()
	if err != nil {
		return err
	}

	var newUsers []models.User
	for _, u := range users {
		if u.ID != id {
			newUsers = append(newUsers, u)
		}
	}

	return userStore.WriteAll(UserList{Users: newUsers})
}
