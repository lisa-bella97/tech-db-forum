package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"net/http"
	"time"
)

func CreatePost(post *models.Post) *models.ModelError {
	dateTimeTemplate := "2006-01-02 15:04:05"
	created := time.Now().Format(dateTimeTemplate)
	err := Connection.QueryRow(`INSERT INTO posts (parent, author, message, forum, thread, created)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created`, post.Parent, post.Author, post.Message, post.Forum,
		post.Thread, created).Scan(&post.Id, &post.Created)
	if err != nil {
		return &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot create post: " + err.Error(),
		}
	}
	return nil
}
