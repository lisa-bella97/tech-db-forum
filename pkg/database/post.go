package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"net/http"
)

func CreatePost(post *models.Post) *models.ModelError {
	err := Connection.QueryRow(`INSERT INTO posts (parent, author, message, forum, thread, created)
		VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id, created`, post.Parent, post.Author, post.Message, post.Forum,
		post.Thread).Scan(&post.Id, &post.Created)
	if err != nil {
		return &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot create post: " + err.Error(),
		}
	}
	return nil
}
