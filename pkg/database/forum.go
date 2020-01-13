package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/pkg/errors"
	"net/http"
)

func GetForumBySlug(slug string) (models.Forum, *models.ModelError) {
	row := Connection.QueryRow(`SELECT * FROM forums WHERE LOWER(slug) = LOWER($1)`, slug)
	forum := models.Forum{}
	err := row.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return models.Forum{}, &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "Can't find forum with slug " + slug,
		}
	}
	return forum, nil
}

func CreateForum(forum models.Forum) error {
	_, err := Connection.Exec(`INSERT INTO forums (title, "user", slug) VALUES ($1, $2, $3)`,
		forum.Title, forum.User, forum.Slug)
	if err != nil {
		return errors.Wrap(err, "cannot create forum")
	}
	return nil
}
