package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/pkg/errors"
)

func GetThreadBySlug(slug string) (models.Thread, error) {
	rows, err := Connection.Query(`SELECT * FROM threads WHERE LOWER(slug) = LOWER($1)`, slug)
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "cannot get thread by slug")
	}
	defer rows.Close()

	if rows.Next() {
		thread := models.Thread{}
		err := rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message,
			&thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return models.Thread{}, errors.Wrap(err, "db query result parsing error")
		}
		return thread, nil
	}

	return models.Thread{}, errors.New("thread not found by slug")
}

func CreateThread(thread *models.Thread) error {
	err := Connection.QueryRow(`INSERT INTO threads (title, author, forum, message, slug, created)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, thread.Title, thread.Author, thread.Forum, thread.Message,
		thread.Slug, thread.Created).Scan(&thread.Id)
	if err != nil {
		return errors.Wrap(err, "cannot create thread")
	}
	return nil
}
