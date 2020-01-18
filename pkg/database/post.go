package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"net/http"
	"time"
)

func CreatePost(post *models.Post) *models.ModelError {
	dateTimeTemplate := "2006-01-02 15:04:05"
	created := time.Now().Format(dateTimeTemplate)
	err := Connection.QueryRow(`INSERT INTO posts (parent, author, message, "isEdited", forum, thread, created, path)
		VALUES ($1, $2, $3, $4, $5, $6, $7, (SELECT path FROM posts WHERE id = $1) ||
		(SELECT last_value FROM posts_id_seq)) RETURNING id, created`, post.Parent, post.Author, post.Message,
		post.IsEdited, post.Forum, post.Thread, created).Scan(&post.Id, &post.Created)
	if err != nil {
		return &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot create post: " + err.Error(),
		}
	}
	return nil
}
