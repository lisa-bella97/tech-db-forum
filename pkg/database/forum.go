package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/pkg/errors"
)

func GetForumBySlug(slug string) (models.Forum, error) {
	res, err := Connection.Query(`SELECT * FROM forums WHERE slug = $1`, slug)
	if err != nil {
		return models.Forum{}, errors.Wrap(err, "cannot get forum by slug")
	}
	defer res.Close()

	f := models.Forum{}

	if res.Next() {
		err := res.Scan(&f.Posts, &f.Slug, &f.Threads, &f.Title, &f.User)
		if err != nil {
			return models.Forum{}, errors.Wrap(err, "db query result parsing error")
		}
	}

	return f, nil
}

func CreateForum(forum models.Forum) error {
	_, err := Connection.Exec(`INSERT INTO forums (posts, slug, threads, title, "user") VALUES ($1, $2, $3, $4, $5)`,
		forum.Posts, forum.Slug, forum.Threads, forum.Title, forum.User)
	if err != nil {
		return errors.Wrap(err, "cannot create forum")
	}

	return nil
}
